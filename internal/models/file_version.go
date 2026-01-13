package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileVersion 文件版本模型
type FileVersion struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FileID        uuid.UUID `gorm:"type:uuid;not null;index" json:"file_id"`
	VersionNumber int       `gorm:"not null" json:"version_number"`
	FileSize      int64     `gorm:"not null" json:"file_size"`
	FileHash      string    `gorm:"type:varchar(64);not null" json:"file_hash"`
	StoragePath   string    `gorm:"type:text;not null" json:"storage_path"`
	MimeType      string    `gorm:"type:varchar(100)" json:"mime_type"`
	ChangeNote    string    `gorm:"type:text" json:"change_note,omitempty"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	File      File     `gorm:"foreignKey:FileID" json:"file,omitempty"`
	Creator   User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName 指定表名
func (FileVersion) TableName() string {
	return "file_versions"
}

// BeforeCreate 创建前的钩子
func (fv *FileVersion) BeforeCreate(tx *gorm.DB) error {
	if fv.ID == uuid.Nil {
		fv.ID = uuid.New()
	}
	return nil
}

// FileVersionCreateRequest 文件版本创建请求
type FileVersionCreateRequest struct {
	ChangeNote string `json:"change_note"`
}

// FileVersionResponse 文件版本响应
type FileVersionResponse struct {
	ID            uuid.UUID `json:"id"`
	FileID        uuid.UUID `json:"file_id"`
	VersionNumber int       `json:"version_number"`
	FileSize      int64     `json:"file_size"`
	FileHash      string    `json:"file_hash"`
	MimeType      string    `json:"mime_type"`
	ChangeNote    string    `json:"change_note,omitempty"`
	CreatedBy     uuid.UUID `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`

	// 可选的关联数据
	FileName    string `json:"file_name,omitempty"`
	CreatorName string `json:"creator_name,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}

// ToResponse 转换为响应格式
func (fv *FileVersion) ToResponse() FileVersionResponse {
	return FileVersionResponse{
		ID:            fv.ID,
		FileID:        fv.FileID,
		VersionNumber: fv.VersionNumber,
		FileSize:      fv.FileSize,
		FileHash:      fv.FileHash,
		MimeType:      fv.MimeType,
		ChangeNote:    fv.ChangeNote,
		CreatedBy:     fv.CreatedBy,
		CreatedAt:     fv.CreatedAt,
	}
}

// GetStoragePath 获取版本文件的存储路径
func (fv *FileVersion) GetStoragePath(basePath string) string {
	return fv.StoragePath
}

// FileVersionFilter 文件版本查询过滤器
type FileVersionFilter struct {
	FileID        *uuid.UUID `form:"file_id"`
	VersionNumber *int       `form:"version_number"`
	CreatedBy     *uuid.UUID `form:"created_by"`
	CreatedAtFrom *time.Time `form:"created_at_from"`
	CreatedAtTo   *time.Time `form:"created_at_to"`
	Page          int        `form:"page" binding:"min=1"`
	PageSize      int        `form:"page_size" binding:"min=1,max=100"`
}

// ApplyFilter 应用过滤器到查询
func (f *FileVersionFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	query := db

	if f.FileID != nil {
		query = query.Where("file_id = ?", *f.FileID)
	}

	if f.VersionNumber != nil {
		query = query.Where("version_number = ?", *f.VersionNumber)
	}

	if f.CreatedBy != nil {
		query = query.Where("created_by = ?", *f.CreatedBy)
	}

	if f.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", *f.CreatedAtFrom)
	}

	if f.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", *f.CreatedAtTo)
	}

	// 默认按版本号降序排序（最新版本在前）
	query = query.Order("version_number DESC")

	return query
}

// VersionInfo 版本信息
type VersionInfo struct {
	CurrentVersion int `json:"current_version"`
	TotalVersions  int `json:"total_versions"`
	MaxVersions    int `json:"max_versions"` // 系统允许的最大版本数
}

// VersionRestoreRequest 版本恢复请求
type VersionRestoreRequest struct {
	VersionNumber int `json:"version_number" binding:"required,min=1"`
}

// VersionCompareResult 版本比较结果
type VersionCompareResult struct {
	Version1 FileVersionResponse `json:"version1"`
	Version2 FileVersionResponse `json:"version2"`
	Differences []VersionDifference `json:"differences"`
}

// VersionDifference 版本差异
type VersionDifference struct {
	Field     string      `json:"field"`
	Value1    interface{} `json:"value1"`
	Value2    interface{} `json:"value2"`
	ChangeType string     `json:"change_type"` // added, modified, removed
}

// FileVersionStats 文件版本统计
type FileVersionStats struct {
	FileID         uuid.UUID `json:"file_id"`
	FileName       string    `json:"file_name"`
	TotalVersions  int       `json:"total_versions"`
	TotalSize      int64     `json:"total_size"`
	FirstVersionAt time.Time `json:"first_version_at"`
	LastVersionAt  time.Time `json:"last_version_at"`
	AverageSize    int64     `json:"average_size"`
}

// CleanupOldVersions 清理旧版本的配置
type CleanupOldVersions struct {
	KeepLastNVersions int `json:"keep_last_n_versions"` // 保留最近N个版本
	MaxAgeDays        int `json:"max_age_days"`         // 最大保留天数
	MinVersions       int `json:"min_versions"`         // 最少保留版本数
}