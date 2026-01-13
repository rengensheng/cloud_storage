package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/models"
)

// FileRepository 文件仓库接口
type FileRepository interface {
	// 基础CRUD操作
	Create(file *models.File) error
	CreateWithTx(tx *gorm.DB, file *models.File) error
	FindByID(id uuid.UUID) (*models.File, error)
	FindByIDIncludingDeleted(id uuid.UUID) (*models.File, error)
	FindAll(filter models.FileFilter) ([]models.File, error)
	FindAllWithTx(tx *gorm.DB, filter models.FileFilter) ([]models.File, error)
	Update(id uuid.UUID, updates map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error
	Delete(id uuid.UUID) error
	DeleteWithTx(tx *gorm.DB, id uuid.UUID) error
	SoftDelete(id uuid.UUID) error
	Restore(id uuid.UUID) error

	// 查询操作
	FindByUserAndName(userID uuid.UUID, parentID *uuid.UUID, name string) (*models.File, error)
	FindByShareToken(token string) (*models.File, error)
	FindOldRecycledFiles(userID uuid.UUID, cutoffDate time.Time) ([]models.File, error)

	// 统计操作
	Count(filter models.FileFilter) (int64, error)
	GetUserFileStats(userID uuid.UUID) (*models.FileStats, error)
}

// fileRepository 文件仓库实现
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库实例
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create 创建文件
func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// CreateWithTx 在事务中创建文件
func (r *fileRepository) CreateWithTx(tx *gorm.DB, file *models.File) error {
	return tx.Create(file).Error
}

// FindByID 根据ID查找文件
func (r *fileRepository) FindByID(id uuid.UUID) (*models.File, error) {
	var file models.File
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByIDIncludingDeleted 根据ID查找文件（包括已删除的）
func (r *fileRepository) FindByIDIncludingDeleted(id uuid.UUID) (*models.File, error) {
	var file models.File
	err := r.db.Unscoped().Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindAll 查找所有符合条件的文件
func (r *fileRepository) FindAll(filter models.FileFilter) ([]models.File, error) {
	var files []models.File

	query := r.db.Model(&models.File{})
	query = filter.ApplyFilter(query)

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return files, nil
}

// FindAllWithTx 在事务中查找所有符合条件的文件
func (r *fileRepository) FindAllWithTx(tx *gorm.DB, filter models.FileFilter) ([]models.File, error) {
	var files []models.File

	query := tx.Model(&models.File{})
	query = filter.ApplyFilter(query)

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Update 更新文件
func (r *fileRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.File{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateWithTx 在事务中更新文件
func (r *fileRepository) UpdateWithTx(tx *gorm.DB, id uuid.UUID, updates map[string]interface{}) error {
	return tx.Model(&models.File{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除文件（硬删除）
func (r *fileRepository) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.File{}, "id = ?", id).Error
}

// DeleteWithTx 在事务中删除文件（硬删除）
func (r *fileRepository) DeleteWithTx(tx *gorm.DB, id uuid.UUID) error {
	return tx.Unscoped().Delete(&models.File{}, "id = ?", id).Error
}

// SoftDelete 软删除文件
func (r *fileRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Delete(&models.File{}, "id = ?", id).Error
}

// Restore 恢复已删除的文件
func (r *fileRepository) Restore(id uuid.UUID) error {
	return r.db.Unscoped().Model(&models.File{}).Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// FindByUserAndName 根据用户ID、父目录ID和文件名查找文件
func (r *fileRepository) FindByUserAndName(userID uuid.UUID, parentID *uuid.UUID, name string) (*models.File, error) {
	var file models.File

	query := r.db.Where("user_id = ? AND name = ? AND deleted_at IS NULL", userID, name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.First(&file).Error
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// FindByShareToken 根据分享令牌查找文件
func (r *fileRepository) FindByShareToken(token string) (*models.File, error) {
	var file models.File
	err := r.db.Where("share_token = ? AND deleted_at IS NULL", token).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindOldRecycledFiles 查找旧的回收站文件
func (r *fileRepository) FindOldRecycledFiles(userID uuid.UUID, cutoffDate time.Time) ([]models.File, error) {
	var files []models.File

	err := r.db.Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL AND deleted_at < ?", userID, cutoffDate).
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}

// Count 统计符合条件的文件数量
func (r *fileRepository) Count(filter models.FileFilter) (int64, error) {
	var count int64

	query := r.db.Model(&models.File{})
	query = filter.ApplyFilter(query)

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetUserFileStats 获取用户文件统计信息
func (r *fileRepository) GetUserFileStats(userID uuid.UUID) (*models.FileStats, error) {
	stats := &models.FileStats{}

	// 统计文件总数
	if err := r.db.Model(&models.File{}).
		Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, models.FileTypeFile).
		Count(&stats.TotalFiles).Error; err != nil {
		return nil, err
	}

	// 统计目录总数
	if err := r.db.Model(&models.File{}).
		Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, models.FileTypeDir).
		Count(&stats.TotalDirs).Error; err != nil {
		return nil, err
	}

	// 统计总大小
	if err := r.db.Model(&models.File{}).
		Select("COALESCE(SUM(size), 0)").
		Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, models.FileTypeFile).
		Scan(&stats.TotalSize).Error; err != nil {
		return nil, err
	}

	// 统计公开文件数
	if err := r.db.Model(&models.File{}).
		Where("user_id = ? AND is_public = ? AND deleted_at IS NULL", userID, true).
		Count(&stats.PublicFiles).Error; err != nil {
		return nil, err
	}

	// 统计最近7天创建的文件数
	weekAgo := time.Now().AddDate(0, 0, -7)
	if err := r.db.Model(&models.File{}).
		Where("user_id = ? AND created_at >= ? AND deleted_at IS NULL", userID, weekAgo).
		Count(&stats.RecentFiles).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// 其他辅助方法

// GetFileTree 获取文件树
func (r *fileRepository) GetFileTree(userID uuid.UUID, rootID *uuid.UUID) ([]*FileTreeNode, error) {
	// 获取根节点文件
	var rootFiles []models.File
	query := r.db.Where("user_id = ? AND deleted_at IS NULL", userID)

	if rootID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", rootID)
	}

	if err := query.Find(&rootFiles).Order("type DESC, name ASC").Error; err != nil {
		return nil, err
	}

	// 构建树节点
	var tree []*FileTreeNode
	for _, file := range rootFiles {
		node := &FileTreeNode{
			File: file,
		}

		// 如果是目录，递归获取子节点
		if file.Type == models.FileTypeDir {
			children, err := r.GetFileTree(userID, &file.ID)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}

		tree = append(tree, node)
	}

	return tree, nil
}

// FileTreeNode 文件树节点
type FileTreeNode struct {
	models.File
	Children []*FileTreeNode `json:"children,omitempty"`
}

// GetFilePath 获取文件完整路径
func (r *fileRepository) GetFilePath(fileID uuid.UUID) (string, error) {
	file, err := r.FindByID(fileID)
	if err != nil {
		return "", err
	}
	return file.Path, nil
}

// GetFileAncestors 获取文件的所有祖先
func (r *fileRepository) GetFileAncestors(fileID uuid.UUID) ([]models.File, error) {
	var ancestors []models.File
	currentID := fileID

	for {
		var file models.File
		err := r.db.Where("id = ?", currentID).First(&file).Error
		if err != nil {
			return nil, err
		}

		// 添加到祖先列表（从父级到根级）
		ancestors = append([]models.File{file}, ancestors...)

		// 如果到达根目录，停止
		if file.ParentID == nil {
			break
		}

		currentID = *file.ParentID
	}

	return ancestors, nil
}

// BulkDelete 批量删除文件
func (r *fileRepository) BulkDelete(fileIDs []uuid.UUID) error {
	return r.db.Where("id IN ?", fileIDs).Delete(&models.File{}).Error
}

// BulkRestore 批量恢复文件
func (r *fileRepository) BulkRestore(fileIDs []uuid.UUID) error {
	return r.db.Unscoped().Model(&models.File{}).
		Where("id IN ?", fileIDs).
		Update("deleted_at", nil).Error
}

// UpdateFileSize 更新文件大小
func (r *fileRepository) UpdateFileSize(fileID uuid.UUID, newSize int64) error {
	return r.db.Model(&models.File{}).
		Where("id = ?", fileID).
		Update("size", newSize).Error
}

// UpdateFileHash 更新文件哈希
func (r *fileRepository) UpdateFileHash(fileID uuid.UUID, hash string) error {
	return r.db.Model(&models.File{}).
		Where("id = ?", fileID).
		Update("hash", hash).Error
}

// GetDuplicateFiles 查找重复文件（根据哈希值）
func (r *fileRepository) GetDuplicateFiles(userID uuid.UUID) (map[string][]models.File, error) {
	// 获取所有有哈希值的文件
	var files []models.File
	err := r.db.Where("user_id = ? AND hash IS NOT NULL AND hash != '' AND deleted_at IS NULL", userID).
		Find(&files).Error
	if err != nil {
		return nil, err
	}

	// 按哈希值分组
	hashMap := make(map[string][]models.File)
	for _, file := range files {
		hashMap[file.Hash] = append(hashMap[file.Hash], file)
	}

	// 过滤出重复的文件
	duplicates := make(map[string][]models.File)
	for hash, fileList := range hashMap {
		if len(fileList) > 1 {
			duplicates[hash] = fileList
		}
	}

	return duplicates, nil
}
