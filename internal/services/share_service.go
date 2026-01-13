package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
	"cloud-storage/internal/repositories"
)

type ShareService struct {
	db        *gorm.DB
	shareRepo repositories.ShareRepository
	fileRepo  repositories.FileRepository
}

func NewShareService(
	db *gorm.DB,
	shareRepo repositories.ShareRepository,
	fileRepo repositories.FileRepository,
) *ShareService {
	return &ShareService{
		db:        db,
		shareRepo: shareRepo,
		fileRepo:  fileRepo,
	}
}

func (s *ShareService) CreateShare(userID uuid.UUID, fileID uuid.UUID, req models.ShareCreateRequest) (*models.Share, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if file.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	var passwordHash *string
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashed := string(hash)
		passwordHash = &hashed
	}

	var expiresAt *time.Time
	if req.ExpiresInDays != nil && *req.ExpiresInDays > 0 {
		expires := time.Now().AddDate(0, 0, *req.ExpiresInDays)
		expiresAt = &expires
	}

	share := &models.Share{
		FileID:       fileID,
		UserID:       userID,
		ShareToken:   generateShareToken(),
		PasswordHash: passwordHash,
		AccessType:   req.AccessType,
		ExpiresAt:    expiresAt,
		MaxDownloads: req.MaxDownloads,
		IsActive:     true,
	}

	if err := s.shareRepo.Create(share); err != nil {
		return nil, fmt.Errorf("failed to create share: %w", err)
	}

	return share, nil
}

func (s *ShareService) GetShare(shareID uuid.UUID, userID uuid.UUID) (*models.Share, error) {
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		return nil, fmt.Errorf("share not found: %w", err)
	}

	if share.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	return share, nil
}

func (s *ShareService) GetUserShares(userID uuid.UUID, filter models.ShareFilter) ([]models.Share, int64, error) {
	filter.Page = 1
	filter.PageSize = 20

	shares, total, err := s.shareRepo.FindByUser(userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get shares: %w", err)
	}

	return shares, total, nil
}

func (s *ShareService) UpdateShare(shareID uuid.UUID, userID uuid.UUID, req models.ShareUpdateRequest) (*models.Share, error) {
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		return nil, fmt.Errorf("share not found: %w", err)
	}

	if share.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	updates := make(map[string]interface{})

	if req.Password != nil {
		if *req.Password == "" {
			updates["password_hash"] = nil
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			updates["password_hash"] = string(hash)
		}
	}

	if req.AccessType != nil {
		updates["access_type"] = *req.AccessType
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if req.ExpiresInDays != nil {
		if *req.ExpiresInDays <= 0 {
			updates["expires_at"] = nil
		} else {
			expires := time.Now().AddDate(0, 0, *req.ExpiresInDays)
			updates["expires_at"] = expires
		}
	}

	if req.MaxDownloads != nil {
		updates["max_downloads"] = *req.MaxDownloads
	}

	if err := s.shareRepo.Update(shareID, updates); err != nil {
		return nil, fmt.Errorf("failed to update share: %w", err)
	}

	return s.shareRepo.FindByID(shareID)
}

func (s *ShareService) DeleteShare(shareID uuid.UUID, userID uuid.UUID) error {
	share, err := s.shareRepo.FindByID(shareID)
	if err != nil {
		return fmt.Errorf("share not found: %w", err)
	}

	if share.UserID != userID {
		return fmt.Errorf("permission denied")
	}

	if err := s.shareRepo.Delete(shareID); err != nil {
		return fmt.Errorf("failed to delete share: %w", err)
	}

	return nil
}

func (s *ShareService) AccessShare(token string, password *string) (*models.Share, error) {
	share, err := s.shareRepo.FindByToken(token)
	if err != nil {
		return nil, fmt.Errorf("share not found")
	}

	if !share.IsValid() {
		return nil, fmt.Errorf("share is invalid or expired")
	}

	if share.PasswordHash != nil {
		if password == nil {
			return nil, fmt.Errorf("password required")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*share.PasswordHash), []byte(*password)); err != nil {
			return nil, fmt.Errorf("invalid password")
		}
	}

	return share, nil
}

func (s *ShareService) DownloadSharedFile(token string, password *string) (*models.File, error) {
	share, err := s.AccessShare(token, password)
	if err != nil {
		return nil, err
	}

	if !share.CanDownload() {
		return nil, fmt.Errorf("download not allowed")
	}

	if err := s.shareRepo.IncrementDownloadCount(share.ID); err != nil {
		return nil, fmt.Errorf("failed to increment download count")
	}

	file, err := s.fileRepo.FindByID(share.FileID)
	if err != nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}

func (s *ShareService) GetShareStats(userID uuid.UUID) (*models.ShareStats, error) {
	stats, err := s.shareRepo.GetUserShareStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get share stats: %w", err)
	}
	return stats, nil
}

func (s *ShareService) BatchDeleteShares(shareIDs []uuid.UUID, userID uuid.UUID) (int, error) {
	deletedCount := 0

	for _, shareID := range shareIDs {
		share, err := s.shareRepo.FindByID(shareID)
		if err != nil {
			continue
		}

		if share.UserID != userID {
			continue
		}

		if err := s.shareRepo.Delete(shareID); err == nil {
			deletedCount++
		}
	}

	return deletedCount, nil
}

func generateShareToken() string {
	token := uuid.New().String()
	token = token[:32]
	return token
}
