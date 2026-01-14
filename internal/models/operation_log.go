package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OperationType 操作类型
type OperationType string

const (
	// 用户相关操作
	OperationUserRegister OperationType = "user_register"
	OperationUserLogin    OperationType = "user_login"
	OperationUserLogout   OperationType = "user_logout"
	OperationUserUpdate   OperationType = "user_update"
	OperationUserDelete   OperationType = "user_delete"

	// 文件相关操作
	OperationFileUpload   OperationType = "file_upload"
	OperationFileDownload OperationType = "download"
	OperationFileCreate   OperationType = "file_create"
	OperationFileUpdate   OperationType = "file_update"
	OperationFileDelete   OperationType = "file_delete"
	OperationFileMove     OperationType = "file_move"
	OperationFileCopy     OperationType = "file_copy"
	OperationFileRename   OperationType = "file_rename"
	OperationFileShare    OperationType = "file_share"
	OperationFileUnshare  OperationType = "file_unshare"

	// 文件夹相关操作
	OperationDirCreate OperationType = "dir_create"
	OperationDirDelete OperationType = "dir_delete"

	// 分享相关操作
	OperationShareCreate OperationType = "share_create"
	OperationShareUpdate OperationType = "share_update"
	OperationShareDelete OperationType = "share_delete"
	OperationShareAccess OperationType = "share_access"

	// 系统操作
	OperationSystemBackup  OperationType = "system_backup"
	OperationSystemRestore OperationType = "system_restore"
	OperationSystemCleanup OperationType = "system_cleanup"
)

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeUser   ResourceType = "user"
	ResourceTypeFile   ResourceType = "file"
	ResourceTypeDir    ResourceType = "directory"
	ResourceTypeShare  ResourceType = "share"
	ResourceTypeSystem ResourceType = "system"
)

// OperationResult 操作结果
type OperationResult string

const (
	OperationSuccess OperationResult = "success"
	OperationFailure OperationResult = "failure"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       *uuid.UUID      `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Operation    OperationType   `gorm:"type:varchar(50);not null;index" json:"operation"`
	ResourceType ResourceType    `gorm:"type:varchar(20);not null" json:"resource_type"`
	ResourceID   *string         `gorm:"type:text" json:"resource_id,omitempty"`
	Result       OperationResult `gorm:"type:varchar(10);not null" json:"result"`
	Details      string          `gorm:"type:text" json:"details,omitempty"`
	IPAddress    string          `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent    string          `gorm:"type:text" json:"user_agent,omitempty"`
	Error        string          `gorm:"type:text" json:"error,omitempty"`
	Duration     int64           `gorm:"default:0" json:"duration"` // 操作耗时，单位毫秒
	CreatedAt    time.Time       `gorm:"autoCreateTime;index" json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// BeforeCreate 创建前的钩子
func (ol *OperationLog) BeforeCreate(tx *gorm.DB) error {
	if ol.ID == uuid.Nil {
		ol.ID = uuid.New()
	}
	return nil
}

// OperationLogCreate 创建操作日志的请求
type OperationLogCreate struct {
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	Operation    OperationType   `json:"operation" binding:"required"`
	ResourceType ResourceType    `json:"resource_type" binding:"required"`
	ResourceID   *string         `json:"resource_id,omitempty"`
	Result       OperationResult `json:"result" binding:"required,oneof=success failure"`
	Details      string          `json:"details,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	Error        string          `json:"error,omitempty"`
	Duration     int64           `json:"duration,omitempty"`
}

// OperationLogResponse 操作日志响应
type OperationLogResponse struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	Operation    OperationType   `json:"operation"`
	ResourceType ResourceType    `json:"resource_type"`
	ResourceID   *string         `json:"resource_id,omitempty"`
	Result       OperationResult `json:"result"`
	Details      string          `json:"details,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	Error        string          `json:"error,omitempty"`
	Duration     int64           `json:"duration"`
	CreatedAt    time.Time       `json:"created_at"`

	// 可选的关联数据
	UserName     string `json:"user_name,omitempty"`
	ResourceName string `json:"resource_name,omitempty"`
}

// ToResponse 转换为响应格式
func (ol *OperationLog) ToResponse() OperationLogResponse {
	return OperationLogResponse{
		ID:           ol.ID,
		UserID:       ol.UserID,
		Operation:    ol.Operation,
		ResourceType: ol.ResourceType,
		ResourceID:   ol.ResourceID,
		Result:       ol.Result,
		Details:      ol.Details,
		IPAddress:    ol.IPAddress,
		UserAgent:    ol.UserAgent,
		Error:        ol.Error,
		Duration:     ol.Duration,
		CreatedAt:    ol.CreatedAt,
	}
}

// OperationLogFilter 操作日志查询过滤器
type OperationLogFilter struct {
	UserID        *uuid.UUID       `form:"-"`
	UserIDStr     string           `form:"user_id"`
	Operation     *OperationType   `form:"operation"`
	ResourceType  *ResourceType    `form:"resource_type"`
	ResourceID    *string          `form:"resource_id"`
	Result        *OperationResult `form:"result"`
	IPAddress     *string          `form:"ip_address"`
	CreatedAtFrom *time.Time       `form:"created_at_from"`
	CreatedAtTo   *time.Time       `form:"created_at_to"`
	Page          int              `form:"page" binding:"omitempty,min=1"`
	PageSize      int              `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy        string           `form:"sort_by" binding:"oneof=created_at operation duration"`
	SortOrder     string           `form:"sort_order" binding:"oneof=asc desc"`
}

// ApplyFilter 应用过滤器到查询
func (f *OperationLogFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	query := db

	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}

	if f.Operation != nil {
		query = query.Where("operation = ?", *f.Operation)
	}

	if f.ResourceType != nil {
		query = query.Where("resource_type = ?", *f.ResourceType)
	}

	if f.ResourceID != nil && *f.ResourceID != "" {
		query = query.Where("resource_id = ?", *f.ResourceID)
	}

	if f.Result != nil {
		query = query.Where("result = ?", *f.Result)
	}

	if f.IPAddress != nil && *f.IPAddress != "" {
		query = query.Where("ip_address = ?", *f.IPAddress)
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
		query = query.Order("created_at DESC") // 默认按时间降序
	}

	return query
}

// OperationStats 操作统计信息
type OperationStats struct {
	TotalOperations int64                   `json:"total_operations"`
	SuccessCount    int64                   `json:"success_count"`
	FailureCount    int64                   `json:"failure_count"`
	ByOperation     map[OperationType]int64 `json:"by_operation"`
	ByResourceType  map[ResourceType]int64  `json:"by_resource_type"`
	ByHour          map[int]int64           `json:"by_hour"` // 24小时分布
	ByDay           map[string]int64        `json:"by_day"`  // 日期分布
}

// AuditLogRequest 审计日志请求
type AuditLogRequest struct {
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	GroupBy   string     `form:"group_by" binding:"oneof=hour day month operation resource_type"`
	Format    string     `form:"format" binding:"oneof=json csv"`
}

// SystemHealthLog 系统健康日志
type SystemHealthLog struct {
	Timestamp       time.Time `json:"timestamp"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsage     float64   `json:"memory_usage"`
	DiskUsage       float64   `json:"disk_usage"`
	ActiveUsers     int       `json:"active_users"`
	ActiveUploads   int       `json:"active_uploads"`
	ActiveDownloads int       `json:"active_downloads"`
	ErrorRate       float64   `json:"error_rate"`
	ResponseTime    float64   `json:"response_time"`
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username  string    `gorm:"type:varchar(50);index" json:"username"`
	IPAddress string    `gorm:"type:varchar(45);index" json:"ip_address"`
	Success   bool      `gorm:"default:false" json:"success"`
	UserAgent string    `gorm:"type:text" json:"user_agent,omitempty"`
	Error     string    `gorm:"type:text" json:"error,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// SecurityAlert 安全警报
type SecurityAlert struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AlertType   string     `gorm:"type:varchar(50);not null" json:"alert_type"`
	Severity    string     `gorm:"type:varchar(20);not null" json:"severity"` // low, medium, high, critical
	Description string     `gorm:"type:text;not null" json:"description"`
	IPAddress   string     `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Details     string     `gorm:"type:jsonb" json:"details,omitempty"`
	Resolved    bool       `gorm:"default:false" json:"resolved"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy  *uuid.UUID `gorm:"type:uuid" json:"resolved_by,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (SecurityAlert) TableName() string {
	return "security_alerts"
}

// SystemStats 系统统计信息
type SystemStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	TotalFiles      int64 `json:"total_files"`
	TotalStorage    int64 `json:"total_storage"`
	TodayOperations int64 `json:"today_operations"`
	ActiveShares    int64 `json:"active_shares"`
}
