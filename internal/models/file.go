package models

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileType 文件类型
type FileType string

const (
	FileTypeFile FileType = "file"
	FileTypeDir  FileType = "directory"
)

// File 文件模型
type File struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Path      string         `gorm:"type:text;not null;index" json:"path"`
	Size      int64          `gorm:"default:0" json:"size"`
	MimeType  string         `gorm:"type:varchar(100)" json:"mime_type"`
	Hash      string         `gorm:"type:varchar(64);index" json:"hash,omitempty"`
	Type      FileType       `gorm:"type:varchar(20);not null" json:"type"`
	IsPublic  bool           `gorm:"default:false" json:"is_public"`
	ShareToken *string       `gorm:"type:varchar(32);uniqueIndex" json:"share_token,omitempty"`
	Version   int            `gorm:"default:1" json:"version"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent      *File          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []File         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Versions    []FileVersion  `gorm:"foreignKey:FileID" json:"versions,omitempty"`
	Shares      []Share        `gorm:"foreignKey:FileID" json:"shares,omitempty"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate 创建前的钩子
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}

	// 设置路径
	if f.Path == "" {
		f.Path = f.BuildPath()
	}

	return nil
}

// BeforeUpdate 更新前的钩子
func (f *File) BeforeUpdate(tx *gorm.DB) error {
	// 如果父目录或名称改变，重新构建路径
	if tx.Statement.Changed("ParentID") || tx.Statement.Changed("Name") {
		f.Path = f.BuildPath()
	}
	return nil
}

// BuildPath 构建文件路径
func (f *File) BuildPath() string {
	if f.ParentID == nil {
		return f.Name
	}

	// 查询父目录路径
	var parentPath string
	if f.Parent != nil && f.Parent.Path != "" {
		parentPath = f.Parent.Path
	} else {
		// 注意：这里简化处理，实际实现中可能需要数据库连接
		// 在GORM的钩子中，可以通过tx参数访问数据库
		// 这里返回相对路径，由上层处理完整路径
		parentPath = ""
	}

	if parentPath == "" {
		return f.Name
	}

	return filepath.Join(parentPath, f.Name)
}

// FileCreateRequest 文件创建请求
type FileCreateRequest struct {
	Name     string    `json:"name" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Type     FileType  `json:"type" binding:"required,oneof=file directory"`
	IsPublic bool      `json:"is_public,omitempty"`
}

// FileUpdateRequest 文件更新请求
type FileUpdateRequest struct {
	Name     *string    `json:"name"`
	ParentID *uuid.UUID `json:"parent_id"`
	IsPublic *bool      `json:"is_public"`
}

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	ParentID *uuid.UUID `form:"parent_id"`
	IsPublic bool       `form:"is_public"`
	Override bool       `form:"override"`
}

// FileResponse 文件响应
type FileResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Size      int64      `json:"size"`
	MimeType  string     `json:"mime_type"`
	Type      FileType   `json:"type"`
	IsPublic  bool       `json:"is_public"`
	ShareToken *string   `json:"share_token,omitempty"`
	Version   int        `json:"version"`
	UserID    uuid.UUID  `json:"user_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// 可选的关联数据
	ChildrenCount int64 `json:"children_count,omitempty"`
	DownloadURL   string `json:"download_url,omitempty"`
	PreviewURL    string `json:"preview_url,omitempty"`
}

// ToResponse 转换为响应格式
func (f *File) ToResponse() FileResponse {
	return FileResponse{
		ID:        f.ID,
		Name:      f.Name,
		Path:      f.Path,
		Size:      f.Size,
		MimeType:  f.MimeType,
		Type:      f.Type,
		IsPublic:  f.IsPublic,
		ShareToken: f.ShareToken,
		Version:   f.Version,
		UserID:    f.UserID,
		ParentID:  f.ParentID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

// IsDirectory 检查是否是目录
func (f *File) IsDirectory() bool {
	return f.Type == FileTypeDir
}

// IsFile 检查是否是文件
func (f *File) IsFile() bool {
	return f.Type == FileTypeFile
}

// GetExtension 获取文件扩展名
func (f *File) GetExtension() string {
	if f.IsDirectory() {
		return ""
	}
	return strings.ToLower(filepath.Ext(f.Name))
}

// GetFullPath 获取完整存储路径
func (f *File) GetFullPath(storagePath string) string {
	return filepath.Join(storagePath, f.UserID.String(), f.Path)
}

// FileFilter 文件查询过滤器
type FileFilter struct {
	UserID    *uuid.UUID `form:"user_id"`
	ParentID  *uuid.UUID `form:"parent_id"`
	Name      *string    `form:"name"`
	Type      *FileType  `form:"type"`
	MimeType  *string    `form:"mime_type"`
	IsPublic  *bool      `form:"is_public"`
	Deleted   *bool      `form:"deleted"`
	CreatedAtFrom *time.Time `form:"created_at_from"`
	CreatedAtTo   *time.Time `form:"created_at_to"`
	Page      int        `form:"page" binding:"min=1"`
	PageSize  int        `form:"page_size" binding:"min=1,max=100"`
	SortBy    string     `form:"sort_by" binding:"oneof=name size created_at updated_at"`
	SortOrder string     `form:"sort_order" binding:"oneof=asc desc"`
}

// ApplyFilter 应用过滤器到查询
func (f *FileFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	query := db

	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}

	if f.ParentID != nil {
		query = query.Where("parent_id = ?", *f.ParentID)
	} else if f.ParentID == nil && f.Deleted != nil && !*f.Deleted {
		// 默认只显示根目录文件（未删除的）
		query = query.Where("parent_id IS NULL")
	}

	if f.Name != nil && *f.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*f.Name+"%")
	}

	if f.Type != nil {
		query = query.Where("type = ?", *f.Type)
	}

	if f.MimeType != nil && *f.MimeType != "" {
		query = query.Where("mime_type ILIKE ?", "%"+*f.MimeType+"%")
	}

	if f.IsPublic != nil {
		query = query.Where("is_public = ?", *f.IsPublic)
	}

	if f.Deleted != nil {
		if *f.Deleted {
			query = query.Unscoped().Where("deleted_at IS NOT NULL")
		} else {
			query = query.Where("deleted_at IS NULL")
		}
	}

	if f.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", *f.CreatedAtFrom)
	}

	if f.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", *f.CreatedAtTo)
	}

	// 排序
	if f.SortBy != "" {
		order := f.SortBy
		if f.SortOrder != "" {
			order = order + " " + f.SortOrder
		}
		query = query.Order(order)
	} else {
		query = query.Order("type DESC, name ASC") // 目录在前，文件在后
	}

	return query
}

// FileStats 文件统计信息
type FileStats struct {
	TotalFiles    int64 `json:"total_files"`
	TotalDirs     int64 `json:"total_dirs"`
	TotalSize     int64 `json:"total_size"`
	PublicFiles   int64 `json:"public_files"`
	RecentFiles   int64 `json:"recent_files"` // 最近7天
}

// FileMoveRequest 文件移动请求
type FileMoveRequest struct {
	TargetParentID *uuid.UUID `json:"target_parent_id" binding:"required"`
}

// FileCopyRequest 文件复制请求
type FileCopyRequest struct {
	TargetParentID *uuid.UUID `json:"target_parent_id" binding:"required"`
	NewName        *string    `json:"new_name"`
}

// FileSearchRequest 文件搜索请求
type FileSearchRequest struct {
	Query    string `form:"q" binding:"required"`
	SearchIn string `form:"search_in" binding:"oneof=name path content"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}