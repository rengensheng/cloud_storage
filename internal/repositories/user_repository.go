package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// 基础CRUD操作
	Create(user *models.User) error
	CreateWithTx(tx *gorm.DB, user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByIDWithTx(tx *gorm.DB, id uuid.UUID) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(filter models.UserFilter) ([]models.User, error)
	Update(id uuid.UUID, updates map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error
	Delete(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error

	// 查询操作
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)

	// 统计操作
	Count(filter models.UserFilter) (int64, error)
	GetUserStats() (*models.UserStats, error)

	// 业务方法
	UpdateLastLogin(id uuid.UUID) error
	UpdateStorageUsage(id uuid.UUID, delta int64) error
	CheckStorageQuota(id uuid.UUID, requiredSize int64) (bool, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// CreateWithTx 在事务中创建用户
func (r *userRepository) CreateWithTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDWithTx 在事务中根据ID查找用户
func (r *userRepository) FindByIDWithTx(tx *gorm.DB, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := tx.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 查找所有符合条件的用户
func (r *userRepository) FindAll(filter models.UserFilter) ([]models.User, error) {
	var users []models.User

	query := r.db.Model(&models.User{})
	query = filter.ApplyFilter(query)

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Update 更新用户
func (r *userRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateWithTx 在事务中更新用户
func (r *userRepository) UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error {
	return tx.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除用户（硬删除）
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.User{}, "id = ?", id).Error
}

// SoftDelete 软删除用户
func (r *userRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("username = ? AND deleted_at IS NULL", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count 统计符合条件的用户数量
func (r *userRepository) Count(filter models.UserFilter) (int64, error) {
	var count int64

	query := r.db.Model(&models.User{})
	query = filter.ApplyFilter(query)

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetUserStats 获取用户统计信息
func (r *userRepository) GetUserStats() (*models.UserStats, error) {
	stats := &models.UserStats{}

	// 统计用户总数
	if err := r.db.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 统计活跃用户数
	if err := r.db.Model(&models.User{}).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// 统计总存储配额
	if err := r.db.Model(&models.User{}).
		Select("COALESCE(SUM(storage_quota), 0)").
		Where("deleted_at IS NULL").
		Scan(&stats.TotalStorage).Error; err != nil {
		return nil, err
	}

	// 统计已使用存储
	if err := r.db.Model(&models.User{}).
		Select("COALESCE(SUM(used_storage), 0)").
		Where("deleted_at IS NULL").
		Scan(&stats.UsedStorage).Error; err != nil {
		return nil, err
	}

	// 计算平均使用率
	if stats.TotalUsers > 0 {
		stats.AverageUsage = stats.UsedStorage / stats.TotalUsers
	}

	return stats, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(id uuid.UUID) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// UpdateStorageUsage 更新存储使用量
func (r *userRepository) UpdateStorageUsage(id uuid.UUID, delta int64) error {
	// 使用SQL表达式确保原子操作
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("used_storage", gorm.Expr("used_storage + ?", delta)).Error
}

// CheckStorageQuota 检查存储配额
func (r *userRepository) CheckStorageQuota(id uuid.UUID, requiredSize int64) (bool, error) {
	var user models.User
	err := r.db.Select("storage_quota", "used_storage").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error
	if err != nil {
		return false, err
	}

	return user.UsedStorage+requiredSize <= user.StorageQuota, nil
}

// 其他辅助方法

// GetUserWithFiles 获取用户及其文件
func (r *userRepository) GetUserWithFiles(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Files", func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL").Order("type DESC, name ASC")
	}).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithStats 获取用户及统计信息
func (r *userRepository) GetUserWithStats(id uuid.UUID) (*models.User, *models.FileStats, error) {
	user, err := r.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	// 获取文件统计（需要文件仓库）
	// 这里简化处理，实际应该调用文件仓库
	stats := &models.FileStats{}

	return user, stats, nil
}

// BulkUpdate 批量更新用户
func (r *userRepository) BulkUpdate(ids []uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).
		Where("id IN ?", ids).
		Updates(updates).Error
}

// BulkDelete 批量删除用户
func (r *userRepository) BulkDelete(ids []uuid.UUID) error {
	return r.db.Where("id IN ?", ids).Delete(&models.User{}).Error
}

// SearchUsers 搜索用户
func (r *userRepository) SearchUsers(query string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User

	// 构建查询条件
	dbQuery := r.db.Where("deleted_at IS NULL").
		Where("username ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%")

	// 获取总数
	var total int64
	if err := dbQuery.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := dbQuery.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetInactiveUsers 获取不活跃用户
func (r *userRepository) GetInactiveUsers(days int) ([]models.User, error) {
	var users []models.User

	// 计算截止日期
	cutoffDate := gorm.Expr("NOW() - INTERVAL '? days'", days)

	err := r.db.Where("deleted_at IS NULL").
		Where("last_login_at IS NULL OR last_login_at < ?", cutoffDate).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserRole 更新用户角色
func (r *userRepository) UpdateUserRole(id uuid.UUID, role models.UserRole) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("role", role).Error
}

// UpdateUserActiveStatus 更新用户活跃状态
func (r *userRepository) UpdateUserActiveStatus(id uuid.UUID, isActive bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

// UpdateUserQuota 更新用户存储配额
func (r *userRepository) UpdateUserQuota(id uuid.UUID, quota int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("storage_quota", quota).Error
}

// GetUsersByRole 根据角色获取用户
func (r *userRepository) GetUsersByRole(role models.UserRole, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User

	// 获取总数
	var total int64
	err := r.db.Model(&models.User{}).
		Where("role = ? AND deleted_at IS NULL", role).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Where("role = ? AND deleted_at IS NULL", role).
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUsersExceedingQuota 获取超出配额的用户
func (r *userRepository) GetUsersExceedingQuota() ([]models.User, error) {
	var users []models.User

	err := r.db.Where("deleted_at IS NULL").
		Where("used_storage > storage_quota").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}