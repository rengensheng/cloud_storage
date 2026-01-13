package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
)

// ShareRepository 分享仓库接口
type ShareRepository interface {
	Create(share *models.Share) error
	FindByID(id uuid.UUID) (*models.Share, error)
	FindByToken(token string) (*models.Share, error)
	FindByUser(userID uuid.UUID, filter models.ShareFilter) ([]models.Share, int64, error)
	FindAll(filter models.ShareFilter) ([]models.Share, error)
	Update(id uuid.UUID, updates map[string]interface{}) error
	Delete(id uuid.UUID) error
	IncrementDownloadCount(id uuid.UUID) error
	GetUserShareStats(userID uuid.UUID) (*models.ShareStats, error)
	FindByFileID(fileID uuid.UUID) ([]models.Share, error)
	UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error
}

type shareRepository struct {
	db *gorm.DB
}

// NewShareRepository 创建分享仓库实例
func NewShareRepository(db *gorm.DB) ShareRepository {
	return &shareRepository{db: db}
}

func (r *shareRepository) Create(share *models.Share) error {
	return r.db.Create(share).Error
}

func (r *shareRepository) FindByID(id uuid.UUID) (*models.Share, error) {
	var share models.Share
	err := r.db.Preload("File").Preload("User").Where("id = ?", id).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *shareRepository) FindByToken(token string) (*models.Share, error) {
	var share models.Share
	err := r.db.Preload("File").Preload("User").Where("share_token = ?", token).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *shareRepository) FindByUser(userID uuid.UUID, filter models.ShareFilter) ([]models.Share, int64, error) {
	var shares []models.Share
	query := r.db.Model(&models.Share{}).Where("user_id = ?", userID)
	query = filter.ApplyFilter(query)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("File").Offset(offset).Limit(filter.PageSize).Find(&shares).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.Model(&models.Share{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return shares, total, nil
}

func (r *shareRepository) FindAll(filter models.ShareFilter) ([]models.Share, error) {
	var shares []models.Share
	query := r.db.Model(&models.Share{})
	query = filter.ApplyFilter(query)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("File").Preload("User").Offset(offset).Limit(filter.PageSize).Find(&shares).Error
	if err != nil {
		return nil, err
	}

	return shares, nil
}

func (r *shareRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.Share{}).Where("id = ?", id).Updates(updates).Error
}

func (r *shareRepository) UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error {
	return tx.Model(&models.Share{}).Where("id = ?", id).Updates(updates).Error
}

func (r *shareRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Share{}, "id = ?", id).Error
}

func (r *shareRepository) IncrementDownloadCount(id uuid.UUID) error {
	return r.db.Model(&models.Share{}).Where("id = ?", id).UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *shareRepository) GetUserShareStats(userID uuid.UUID) (*models.ShareStats, error) {
	stats := &models.ShareStats{}

	if err := r.db.Model(&models.Share{}).Where("user_id = ?", userID).Count(&stats.TotalShares).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	if err := r.db.Model(&models.Share{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("(expires_at IS NULL OR expires_at >= ?)", now).
		Count(&stats.ActiveShares).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Share{}).
		Where("user_id = ? AND expires_at < ?", userID, now).
		Count(&stats.ExpiredShares).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Share{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(download_count), 0)").
		Scan(&stats.TotalDownloads).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.Share{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&stats.PublicFiles).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *shareRepository) FindByFileID(fileID uuid.UUID) ([]models.Share, error) {
	var shares []models.Share
	err := r.db.Where("file_id = ?", fileID).Find(&shares).Error
	if err != nil {
		return nil, err
	}
	return shares, nil
}
