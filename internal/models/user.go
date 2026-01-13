package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole 用户角色类型
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User 用户模型
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole       `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	StorageQuota int64          `gorm:"default:10737418240" json:"storage_quota"` // 10GB默认
	UsedStorage  int64          `gorm:"default:0" json:"used_storage"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Files       []File       `gorm:"foreignKey:UserID" json:"files,omitempty"`
	Shares      []Share      `gorm:"foreignKey:UserID" json:"shares,omitempty"`
	OperationLogs []OperationLog `gorm:"foreignKey:UserID" json:"operation_logs,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     UserRole `json:"role"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Username     *string  `json:"username"`
	Email        *string  `json:"email"`
	Password     *string  `json:"password"`
	Role         *UserRole `json:"role"`
	StorageQuota *int64   `json:"storage_quota"`
	IsActive     *bool    `json:"is_active"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Role         UserRole   `json:"role"`
	StorageQuota int64      `json:"storage_quota"`
	UsedStorage  int64      `json:"used_storage"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Role:         u.Role,
		StorageQuota: u.StorageQuota,
		UsedStorage:  u.UsedStorage,
		IsActive:     u.IsActive,
		LastLoginAt:  u.LastLoginAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// CheckStorageQuota 检查存储配额
func (u *User) CheckStorageQuota(fileSize int64) bool {
	return u.UsedStorage+fileSize <= u.StorageQuota
}

// UpdateUsedStorage 更新已使用存储空间
func (u *User) UpdateUsedStorage(tx *gorm.DB, delta int64) error {
	u.UsedStorage += delta
	if u.UsedStorage < 0 {
		u.UsedStorage = 0
	}
	return tx.Save(u).Error
}

// HasPermission 检查用户权限
func (u *User) HasPermission(requiredRole UserRole) bool {
	// 权限层级：admin > user
	roleHierarchy := map[UserRole]int{
		RoleUser:  1,
		RoleAdmin: 2,
	}

	return roleHierarchy[u.Role] >= roleHierarchy[requiredRole]
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UserStats 用户统计信息
type UserStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	TotalStorage    int64 `json:"total_storage"`
	UsedStorage     int64 `json:"used_storage"`
	AverageUsage    int64 `json:"average_usage"`
}

// UserFilter 用户查询过滤器
type UserFilter struct {
	Username  *string  `form:"username"`
	Email     *string  `form:"email"`
	Role      *UserRole `form:"role"`
	IsActive  *bool    `form:"is_active"`
	CreatedAtFrom *time.Time `form:"created_at_from"`
	CreatedAtTo   *time.Time `form:"created_at_to"`
	Page      int      `form:"page" binding:"min=1"`
	PageSize  int      `form:"page_size" binding:"min=1,max=100"`
}

// ApplyFilter 应用过滤器到查询
func (f *UserFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	query := db

	if f.Username != nil && *f.Username != "" {
		query = query.Where("username ILIKE ?", "%"+*f.Username+"%")
	}

	if f.Email != nil && *f.Email != "" {
		query = query.Where("email ILIKE ?", "%"+*f.Email+"%")
	}

	if f.Role != nil {
		query = query.Where("role = ?", *f.Role)
	}

	if f.IsActive != nil {
		query = query.Where("is_active = ?", *f.IsActive)
	}

	if f.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", *f.CreatedAtFrom)
	}

	if f.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", *f.CreatedAtTo)
	}

	return query
}