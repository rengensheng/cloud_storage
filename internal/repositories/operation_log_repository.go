package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
)

type OperationLogRepository interface {
	Create(log *models.OperationLog) error
	FindByID(id uuid.UUID) (*models.OperationLog, error)
	FindByUser(userID uuid.UUID, filter models.OperationLogFilter) ([]models.OperationLog, int64, error)
	FindAll(filter models.OperationLogFilter) ([]models.OperationLog, int64, error)
	Delete(id uuid.UUID) error
	DeleteOldLogs(beforeDate time.Time) (int64, error)
	GetUserOperationStats(userID uuid.UUID, startDate, endDate time.Time) (map[string]int64, error)
	GetSystemStats() (*models.SystemStats, error)
}

type operationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(log *models.OperationLog) error {
	return r.db.Create(log).Error
}

func (r *operationLogRepository) FindByID(id uuid.UUID) (*models.OperationLog, error) {
	var log models.OperationLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *operationLogRepository) FindByUser(userID uuid.UUID, filter models.OperationLogFilter) ([]models.OperationLog, int64, error) {
	var logs []models.OperationLog
	query := r.db.Model(&models.OperationLog{}).Where("user_id = ?", userID)
	query = filter.ApplyFilter(query)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.Model(&models.OperationLog{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *operationLogRepository) FindAll(filter models.OperationLogFilter) ([]models.OperationLog, int64, error) {
	var logs []models.OperationLog
	query := r.db.Model(&models.OperationLog{})
	query = filter.ApplyFilter(query)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.Model(&models.OperationLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *operationLogRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.OperationLog{}, "id = ?", id).Error
}

func (r *operationLogRepository) DeleteOldLogs(beforeDate time.Time) (int64, error) {
	result := r.db.Where("created_at < ?", beforeDate).Delete(&models.OperationLog{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *operationLogRepository) GetUserOperationStats(userID uuid.UUID, startDate, endDate time.Time) (map[string]int64, error) {
	type OperationCount struct {
		OperationType string `gorm:"column:operation_type"`
		Count         int64  `gorm:"column:count"`
	}

	var counts []OperationCount
	err := r.db.Model(&models.OperationLog{}).
		Select("operation_type, COUNT(*) as count").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("operation_type").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, c := range counts {
		stats[c.OperationType] = c.Count
	}

	return stats, nil
}

func (r *operationLogRepository) GetSystemStats() (*models.SystemStats, error) {
	stats := &models.SystemStats{}

	// 统计总用户数
	if err := r.db.Table("users").Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 统计活跃用户数
	if err := r.db.Table("users").Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// 统计总文件数
	if err := r.db.Table("files").Where("deleted_at IS NULL").Count(&stats.TotalFiles).Error; err != nil {
		return nil, err
	}

	// 统计总文件大小
	if err := r.db.Table("files").
		Select("COALESCE(SUM(size), 0)").
		Where("deleted_at IS NULL").
		Scan(&stats.TotalStorage).Error; err != nil {
		return nil, err
	}

	// 统计今日操作数
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.Table("operation_logs").
		Where("created_at >= ?", today).
		Count(&stats.TodayOperations).Error; err != nil {
		return nil, err
	}

	// 统计总分享数
	if err := r.db.Table("shares").Where("is_active = ?", true).Count(&stats.ActiveShares).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
