package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
)

type FileVersionRepository interface {
	Create(version *models.FileVersion) error
	FindByID(id uuid.UUID) (*models.FileVersion, error)
	FindByFileID(fileID uuid.UUID) ([]models.FileVersion, error)
	FindByVersion(fileID uuid.UUID, versionNumber int) (*models.FileVersion, error)
	Delete(id uuid.UUID) error
	DeleteByFileID(fileID uuid.UUID) error
}

type fileVersionRepository struct {
	db *gorm.DB
}

func NewFileVersionRepository(db *gorm.DB) FileVersionRepository {
	return &fileVersionRepository{db: db}
}

func (r *fileVersionRepository) Create(version *models.FileVersion) error {
	return r.db.Create(version).Error
}

func (r *fileVersionRepository) FindByID(id uuid.UUID) (*models.FileVersion, error) {
	var version models.FileVersion
	err := r.db.Where("id = ?", id).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *fileVersionRepository) FindByFileID(fileID uuid.UUID) ([]models.FileVersion, error) {
	var versions []models.FileVersion
	err := r.db.Where("file_id = ?", fileID).Order("version_number DESC").Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (r *fileVersionRepository) FindByVersion(fileID uuid.UUID, versionNumber int) (*models.FileVersion, error) {
	var version models.FileVersion
	err := r.db.Where("file_id = ? AND version_number = ?", fileID, versionNumber).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *fileVersionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.FileVersion{}, "id = ?", id).Error
}

func (r *fileVersionRepository) DeleteByFileID(fileID uuid.UUID) error {
	return r.db.Where("file_id = ?", fileID).Delete(&models.FileVersion{}).Error
}
