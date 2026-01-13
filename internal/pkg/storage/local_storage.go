package storage

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	config StorageConfig
	mu     sync.RWMutex
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(config StorageConfig) (*LocalStorage, error) {
	// 确保存储目录存在
	if err := os.MkdirAll(config.LocalPath, 0755); err != nil {
		return nil, wrapStorageError("failed to create storage directory", err)
	}

	return &LocalStorage{
		config: config,
	}, nil
}

// Type 返回存储类型
func (s *LocalStorage) Type() StorageType {
	return StorageTypeLocal
}

// Config 返回存储配置
func (s *LocalStorage) Config() StorageConfig {
	return s.config
}

// Save 保存文件
func (s *LocalStorage) Save(ctx context.Context, key string, data io.Reader, size int64) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}

	filePath := s.getFilePath(key)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return wrapStorageError("failed to create directory", err)
	}

	// 创建临时文件
	tempFile := filePath + ".tmp"
	f, err := os.Create(tempFile)
	if err != nil {
		return wrapStorageError("failed to create temp file", err)
	}
	defer f.Close()

	// 写入数据
	if _, err := io.Copy(f, data); err != nil {
		os.Remove(tempFile)
		return wrapStorageError("failed to write file", err)
	}

	// 关闭文件
	if err := f.Close(); err != nil {
		os.Remove(tempFile)
		return wrapStorageError("failed to close file", err)
	}

	// 重命名为最终文件
	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return wrapStorageError("failed to rename file", err)
	}

	return nil
}

// Get 获取文件
func (s *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	filePath := s.getFilePath(key)

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, wrapStorageError("failed to open file", err)
	}

	return f, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}

	filePath := s.getFilePath(key)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return wrapStorageError("failed to delete file", err)
	}

	// 尝试删除空目录
	s.cleanupEmptyDirs(filepath.Dir(filePath))

	return nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	if !IsValidKey(key) {
		return false, ErrInvalidKey
	}

	filePath := s.getFilePath(key)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, wrapStorageError("failed to stat file", err)
	}

	return true, nil
}

// Stat 获取文件信息
func (s *LocalStorage) Stat(ctx context.Context, key string) (*FileInfo, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	filePath := s.getFilePath(key)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, wrapStorageError("failed to stat file", err)
	}

	// 计算ETag
	etag, err := s.calculateETag(filePath)
	if err != nil {
		etag = fmt.Sprintf("%x", md5.Sum([]byte(info.Name()+info.ModTime().String())))
	}

	return &FileInfo{
		Path:         key,
		Size:         info.Size(),
		LastModified: info.ModTime().Unix(),
		IsDir:        info.IsDir(),
		MimeType:     GetMimeType(filepath.Base(key)),
		ETag:         etag,
	}, nil
}

// Copy 复制文件
func (s *LocalStorage) Copy(ctx context.Context, srcKey, dstKey string) error {
	if !IsValidKey(srcKey) || !IsValidKey(dstKey) {
		return ErrInvalidKey
	}

	srcPath := s.getFilePath(srcKey)
	dstPath := s.getFilePath(dstKey)

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return wrapStorageError("failed to create directory", err)
	}

	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return wrapStorageError("failed to open source file", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return wrapStorageError("failed to create destination file", err)
	}
	defer dstFile.Close()

	// 复制数据
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		os.Remove(dstPath)
		return wrapStorageError("failed to copy file", err)
	}

	return nil
}

// Move 移动文件
func (s *LocalStorage) Move(ctx context.Context, srcKey, dstKey string) error {
	if !IsValidKey(srcKey) || !IsValidKey(dstKey) {
		return ErrInvalidKey
	}

	srcPath := s.getFilePath(srcKey)
	dstPath := s.getFilePath(dstKey)

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return wrapStorageError("failed to create directory", err)
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return wrapStorageError("failed to move file", err)
	}

	// 清理源目录
	s.cleanupEmptyDirs(filepath.Dir(srcPath))

	return nil
}

// List 列出文件
func (s *LocalStorage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	if prefix != "" && !IsValidKey(prefix) {
		return nil, ErrInvalidKey
	}

	dirPath := s.getFilePath(prefix)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []FileInfo{}, nil
		}
		return nil, wrapStorageError("failed to list directory", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		entryPath := filepath.Join(prefix, entry.Name())
		fullPath := filepath.Join(dirPath, entry.Name())

		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		// 计算ETag
		etag, _ := s.calculateETag(fullPath)
		if etag == "" {
			etag = fmt.Sprintf("%x", md5.Sum([]byte(entry.Name()+info.ModTime().String())))
		}

		files = append(files, FileInfo{
			Path:         entryPath,
			Size:         info.Size(),
			LastModified: info.ModTime().Unix(),
			IsDir:        entry.IsDir(),
			MimeType:     GetMimeType(entry.Name()),
			ETag:         etag,
		})
	}

	return files, nil
}

// CreateDir 创建目录
func (s *LocalStorage) CreateDir(ctx context.Context, path string) error {
	if !IsValidKey(path) {
		return ErrInvalidKey
	}

	dirPath := s.getFilePath(path)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return wrapStorageError("failed to create directory", err)
	}

	return nil
}

// DeleteDir 删除目录
func (s *LocalStorage) DeleteDir(ctx context.Context, path string) error {
	if !IsValidKey(path) {
		return ErrInvalidKey
	}

	dirPath := s.getFilePath(path)

	if err := os.RemoveAll(dirPath); err != nil {
		return wrapStorageError("failed to delete directory", err)
	}

	// 清理父目录
	s.cleanupEmptyDirs(filepath.Dir(dirPath))

	return nil
}

// InitiateMultipartUpload 初始化分片上传（本地存储简化实现）
func (s *LocalStorage) InitiateMultipartUpload(ctx context.Context, key string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	uploadID := uuid.New().String()
	tempDir := s.getMultipartUploadDir(uploadID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", wrapStorageError("failed to create multipart upload directory", err)
	}

	return uploadID, nil
}

// UploadPart 上传分片
func (s *LocalStorage) UploadPart(ctx context.Context, uploadID string, partNumber int, data io.Reader) (string, error) {
	tempDir := s.getMultipartUploadDir(uploadID)

	// 检查上传是否存在
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return "", wrapStorageError("multipart upload not found", err)
	}

	// 保存分片
	partFile := filepath.Join(tempDir, fmt.Sprintf("part-%d", partNumber))
	f, err := os.Create(partFile)
	if err != nil {
		return "", wrapStorageError("failed to create part file", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		os.Remove(partFile)
		return "", wrapStorageError("failed to write part file", err)
	}

	// 计算分片ETag（使用MD5）
	if _, err := f.Seek(0, 0); err != nil {
		return "", wrapStorageError("failed to seek part file", err)
	}

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", wrapStorageError("failed to calculate hash", err)
	}

	etag := fmt.Sprintf("%x", hash.Sum(nil))
	return etag, nil
}

// CompleteMultipartUpload 完成分片上传
func (s *LocalStorage) CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) error {
	tempDir := s.getMultipartUploadDir(uploadID)
	keyFile := filepath.Join(tempDir, "key.txt")

	// 读取目标键
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return wrapStorageError("failed to read upload key", err)
	}
	key := string(keyBytes)

	// 合并分片
	filePath := s.getFilePath(key)
	tempFilePath := filePath + ".tmp"

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return wrapStorageError("failed to create directory", err)
	}

	// 创建目标文件
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		return wrapStorageError("failed to create output file", err)
	}
	defer outFile.Close()

	// 按顺序合并分片
	for i := 1; i <= len(parts); i++ {
		partFile := filepath.Join(tempDir, fmt.Sprintf("part-%d", i))

		inFile, err := os.Open(partFile)
		if err != nil {
			return wrapStorageError(fmt.Sprintf("failed to open part %d", i), err)
		}

		if _, err := io.Copy(outFile, inFile); err != nil {
			inFile.Close()
			return wrapStorageError(fmt.Sprintf("failed to copy part %d", i), err)
		}

		inFile.Close()
	}

	// 关闭并重命名
	if err := outFile.Close(); err != nil {
		return wrapStorageError("failed to close output file", err)
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		return wrapStorageError("failed to rename output file", err)
	}

	// 清理临时目录
	os.RemoveAll(tempDir)

	return nil
}

// AbortMultipartUpload 中止分片上传
func (s *LocalStorage) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	tempDir := s.getMultipartUploadDir(uploadID)

	if err := os.RemoveAll(tempDir); err != nil {
		return wrapStorageError("failed to abort multipart upload", err)
	}

	return nil
}

// GetURL 获取文件URL（本地存储返回文件路径）
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	return s.getFilePath(key), nil
}

// GetDownloadURL 获取下载URL
func (s *LocalStorage) GetDownloadURL(ctx context.Context, key string, filename string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	// 本地存储直接返回文件路径
	return s.getFilePath(key), nil
}

// 辅助方法

// getFilePath 获取文件的完整路径
func (s *LocalStorage) getFilePath(key string) string {
	return filepath.Join(s.config.LocalPath, key)
}

// getMultipartUploadDir 获取分片上传临时目录
func (s *LocalStorage) getMultipartUploadDir(uploadID string) string {
	return filepath.Join(s.config.LocalPath, ".multipart", uploadID)
}

// calculateETag 计算文件ETag
func (s *LocalStorage) calculateETag(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// cleanupEmptyDirs 清理空目录
func (s *LocalStorage) cleanupEmptyDirs(dir string) {
	for {
		// 检查是否在存储根目录内
		rel, err := filepath.Rel(s.config.LocalPath, dir)
		if err != nil || rel == "." || rel == ".." || filepath.HasPrefix(rel, "..") {
			break
		}

		// 尝试列出目录内容
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}

		// 删除空目录
		if err := os.Remove(dir); err != nil {
			break
		}

		// 向上级目录继续清理
		dir = filepath.Dir(dir)
	}
}

// DiskUsage 获取磁盘使用情况
type DiskUsage struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
}

// GetDiskUsage 获取磁盘使用情况
func (s *LocalStorage) GetDiskUsage() (*DiskUsage, error) {
	// 这里简化实现，实际应该调用系统API
	var stat syscall.Statfs_t
	err := syscall.Statfs(s.config.LocalPath, &stat)
	if err != nil {
		return nil, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free

	return &DiskUsage{
		Total: total,
		Used:  used,
		Free:  free,
	}, nil
}
