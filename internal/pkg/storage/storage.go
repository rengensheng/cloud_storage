package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeMinIO StorageType = "minio"
)

// StorageConfig 存储配置
type StorageConfig struct {
	Type       StorageType
	LocalPath  string
	Bucket     string
	Region     string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string
	Size         int64
	LastModified int64
	IsDir        bool
	MimeType     string
	ETag         string
}

// Storage 存储接口
type Storage interface {
	// 基础操作
	Type() StorageType
	Config() StorageConfig

	// 文件操作
	Save(ctx context.Context, key string, data io.Reader, size int64) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Stat(ctx context.Context, key string) (*FileInfo, error)
	Copy(ctx context.Context, srcKey, dstKey string) error
	Move(ctx context.Context, srcKey, dstKey string) error

	// 目录操作
	List(ctx context.Context, prefix string) ([]FileInfo, error)
	CreateDir(ctx context.Context, path string) error
	DeleteDir(ctx context.Context, path string) error

	// 分片上传
	InitiateMultipartUpload(ctx context.Context, key string) (string, error)
	UploadPart(ctx context.Context, uploadID string, partNumber int, data io.Reader) (string, error)
	CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) error
	AbortMultipartUpload(ctx context.Context, uploadID string) error

	// 工具方法
	GetURL(ctx context.Context, key string) (string, error)
	GetDownloadURL(ctx context.Context, key string, filename string) (string, error)
}

// NewStorage 创建存储实例
func NewStorage(config StorageConfig) (Storage, error) {
	switch config.Type {
	case StorageTypeLocal:
		return NewLocalStorage(config)
	case StorageTypeS3:
		return NewS3Storage(config)
	case StorageTypeMinIO:
		return NewMinIOStorage(config)
	default:
		return nil, ErrUnsupportedStorageType
	}
}

// 错误定义
var (
	ErrUnsupportedStorageType = newStorageError("unsupported storage type")
	ErrFileNotFound           = newStorageError("file not found")
	ErrPermissionDenied       = newStorageError("permission denied")
	ErrStorageFull            = newStorageError("storage is full")
	ErrInvalidKey             = newStorageError("invalid key")
	ErrUploadFailed           = newStorageError("upload failed")
	ErrDownloadFailed         = newStorageError("download failed")
	ErrDeleteFailed           = newStorageError("delete failed")
)

// storageError 存储错误
type storageError struct {
	message string
	cause   error
}

func newStorageError(message string) *storageError {
	return &storageError{message: message}
}

func wrapStorageError(message string, cause error) *storageError {
	return &storageError{message: message, cause: cause}
}

func (e *storageError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *storageError) Unwrap() error {
	return e.cause
}

// IsStorageError 检查错误是否是存储错误
func IsStorageError(err error) bool {
	_, ok := err.(*storageError)
	return ok
}

// 工具函数

// GenerateFileKey 生成文件存储键
func GenerateFileKey(userID uuid.UUID, filePath string) string {
	return filepath.Join(userID.String(), filePath)
}

// GenerateTempKey 生成临时文件键
func GenerateTempKey(userID uuid.UUID, filename string) string {
	tempID := uuid.New().String()
	return filepath.Join("temp", userID.String(), tempID, filename)
}

// GenerateVersionKey 生成版本文件键
func GenerateVersionKey(userID uuid.UUID, fileID uuid.UUID, version int) string {
	return filepath.Join("versions", userID.String(), fileID.String(),
		fmt.Sprintf("v%d", version))
}

// EnsureDir 确保目录存在
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// CleanPath 清理路径
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// IsValidKey 检查键是否有效
func IsValidKey(key string) bool {
	// 防止路径遍历攻击
	if filepath.IsAbs(key) {
		return false
	}

	cleanPath := filepath.Clean(key)
	if cleanPath != key {
		return false
	}

	// 检查是否包含 ".."
	if strings.HasPrefix(key, "..") || strings.Contains(key, "..") {
		return false
	}

	return true
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// GetMimeType 根据扩展名获取MIME类型
func GetMimeType(filename string) string {
	ext := GetFileExtension(filename)
	switch ext {
	case ".txt", ".md":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".zip":
		return "application/zip"
	case ".tar":
		return "application/x-tar"
	case ".gz":
		return "application/gzip"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	case ".ppt", ".pptx":
		return "application/vnd.ms-powerpoint"
	default:
		return "application/octet-stream"
	}
}