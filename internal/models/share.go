package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShareAccessType 分享访问类型
type ShareAccessType string

const (
	ShareAccessView     ShareAccessType = "view"
	ShareAccessDownload ShareAccessType = "download"
	ShareAccessEdit     ShareAccessType = "edit"
)

// Share 分享模型
type Share struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FileID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"file_id"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	ShareToken    string          `gorm:"type:varchar(32);uniqueIndex;not null" json:"share_token"`
	PasswordHash  *string         `gorm:"type:varchar(255)" json:"-"`
	AccessType    ShareAccessType `gorm:"type:varchar(20);default:'view'" json:"access_type"`
	ExpiresAt     *time.Time      `gorm:"index" json:"expires_at,omitempty"`
	MaxDownloads  *int            `json:"max_downloads,omitempty"`
	DownloadCount int             `gorm:"default:0" json:"download_count"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Share) TableName() string {
	return "shares"
}

// BeforeCreate 创建前的钩子
func (s *Share) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// ShareCreateRequest 分享创建请求
type ShareCreateRequest struct {
	FileID        uuid.UUID       `json:"file_id" binding:"required"`
	Password      *string         `json:"password,omitempty"`
	AccessType    ShareAccessType `json:"access_type" binding:"oneof=view download edit"`
	ExpiresInDays *int            `json:"expires_in_days,omitempty" binding:"omitempty,min=1,max=365"`
	MaxDownloads  *int            `json:"max_downloads,omitempty" binding:"omitempty,min=1"`
}

// ShareUpdateRequest 分享更新请求
type ShareUpdateRequest struct {
	Password      *string          `json:"password,omitempty"`
	AccessType    *ShareAccessType `json:"access_type"`
	IsActive      *bool            `json:"is_active"`
	ExpiresInDays *int             `json:"expires_in_days,omitempty" binding:"omitempty,min=1,max=365"`
	MaxDownloads  *int             `json:"max_downloads,omitempty" binding:"omitempty,min=1"`
}

// ShareResponse 分享响应
type ShareResponse struct {
	ID            uuid.UUID       `json:"id"`
	FileID        uuid.UUID       `json:"file_id"`
	UserID        uuid.UUID       `json:"user_id"`
	ShareToken    string          `json:"share_token"`
	AccessType    ShareAccessType `json:"access_type"`
	ExpiresAt     *time.Time      `json:"expires_at,omitempty"`
	MaxDownloads  *int            `json:"max_downloads,omitempty"`
	DownloadCount int             `json:"download_count"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`

	// 可选的关联数据
	FileName           string `json:"file_name,omitempty"`
	FileSize           int64  `json:"file_size,omitempty"`
	FileType           string `json:"file_type,omitempty"`
	UserName           string `json:"user_name,omitempty"`
	ShareURL           string `json:"share_url,omitempty"`
	HasPassword        bool   `json:"has_password"`
	IsExpired          bool   `json:"is_expired"`
	RemainingDownloads *int   `json:"remaining_downloads,omitempty"`
}

// ToResponse 转换为响应格式
func (s *Share) ToResponse() ShareResponse {
	hasPassword := s.PasswordHash != nil && *s.PasswordHash != ""

	// 检查是否过期
	isExpired := false
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		isExpired = true
	}

	// 计算剩余下载次数
	var remainingDownloads *int
	if s.MaxDownloads != nil {
		remaining := *s.MaxDownloads - s.DownloadCount
		if remaining < 0 {
			remaining = 0
		}
		remainingDownloads = &remaining
	}

	return ShareResponse{
		ID:                 s.ID,
		FileID:             s.FileID,
		UserID:             s.UserID,
		ShareToken:         s.ShareToken,
		AccessType:         s.AccessType,
		ExpiresAt:          s.ExpiresAt,
		MaxDownloads:       s.MaxDownloads,
		DownloadCount:      s.DownloadCount,
		IsActive:           s.IsActive,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		HasPassword:        hasPassword,
		IsExpired:          isExpired,
		RemainingDownloads: remainingDownloads,
	}
}

// IsValid 检查分享是否有效
func (s *Share) IsValid() bool {
	if !s.IsActive {
		return false
	}

	// 检查过期时间
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		return false
	}

	// 检查下载次数限制
	if s.MaxDownloads != nil && s.DownloadCount >= *s.MaxDownloads {
		return false
	}

	return true
}

// CanDownload 检查是否可以下载
func (s *Share) CanDownload() bool {
	if !s.IsValid() {
		return false
	}
	return s.AccessType == ShareAccessDownload || s.AccessType == ShareAccessEdit
}

// CanEdit 检查是否可以编辑
func (s *Share) CanEdit() bool {
	if !s.IsValid() {
		return false
	}
	return s.AccessType == ShareAccessEdit
}

// IncrementDownloadCount 增加下载计数
func (s *Share) IncrementDownloadCount() error {
	s.DownloadCount++
	return nil
}

// ShareFilter 分享查询过滤器
type ShareFilter struct {
	UserID        *uuid.UUID       `form:"-"`
	FileID        *uuid.UUID       `form:"-"`
	UserIDStr     string           `form:"user_id"`
	FileIDStr     string           `form:"file_id"`
	AccessType    *ShareAccessType `form:"access_type"`
	IsActive      *bool            `form:"is_active"`
	Expired       *bool            `form:"expired"`
	CreatedAtFrom *time.Time       `form:"created_at_from"`
	CreatedAtTo   *time.Time       `form:"created_at_to"`
	Page          int              `form:"page" binding:"omitempty,min=1"`
	PageSize      int              `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// ApplyFilter 应用过滤器到查询
func (f *ShareFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	query := db

	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}

	if f.FileID != nil {
		query = query.Where("file_id = ?", *f.FileID)
	}

	if f.AccessType != nil {
		query = query.Where("access_type = ?", *f.AccessType)
	}

	if f.IsActive != nil {
		query = query.Where("is_active = ?", *f.IsActive)
	}

	if f.Expired != nil {
		now := time.Now()
		if *f.Expired {
			query = query.Where("expires_at < ?", now)
		} else {
			query = query.Where("expires_at IS NULL OR expires_at >= ?", now)
		}
	}

	if f.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", *f.CreatedAtFrom)
	}

	if f.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", *f.CreatedAtTo)
	}

	// 默认按创建时间降序排序
	query = query.Order("created_at DESC")

	return query
}

// ShareStats 分享统计信息
type ShareStats struct {
	TotalShares    int64 `json:"total_shares"`
	ActiveShares   int64 `json:"active_shares"`
	ExpiredShares  int64 `json:"expired_shares"`
	TotalDownloads int64 `json:"total_downloads"`
	PublicFiles    int64 `json:"public_files"` // 通过分享可访问的文件
}

// ShareAccessRequest 分享访问请求
type ShareAccessRequest struct {
	Token    string  `json:"token" binding:"required"`
	Password *string `json:"password,omitempty"`
}

// ShareAccessResponse 分享访问响应
type ShareAccessResponse struct {
	Share       ShareResponse `json:"share"`
	File        FileResponse  `json:"file"`
	AccessURL   string        `json:"access_url"`
	CanDownload bool          `json:"can_download"`
	CanEdit     bool          `json:"can_edit"`
	ExpiresIn   *string       `json:"expires_in,omitempty"` // 剩余时间，如 "3天"
}

// ShareLinkInfo 分享链接信息
type ShareLinkInfo struct {
	Token    string `json:"token"`
	URL      string `json:"url"`
	QRCode   string `json:"qr_code,omitempty"` // Base64编码的QR码
	ShortURL string `json:"short_url,omitempty"`
}

// ShareBulkDeleteRequest 批量删除分享请求
type ShareBulkDeleteRequest struct {
	ShareIDs []uuid.UUID `json:"share_ids" binding:"required,min=1"`
}
