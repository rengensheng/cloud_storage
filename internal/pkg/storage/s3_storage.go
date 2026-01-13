package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Storage S3存储实现
type S3Storage struct {
	config StorageConfig
	client *s3.S3
}

// NewS3Storage 创建S3存储实例
func NewS3Storage(config StorageConfig) (*S3Storage, error) {
	// 创建AWS会话
	awsConfig := &aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKey,
			config.SecretKey,
			"",
		),
	}

	// 如果提供了自定义端点（用于MinIO等）
	if config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
		awsConfig.DisableSSL = aws.Bool(!config.UseSSL)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, wrapStorageError("failed to create AWS session", err)
	}

	client := s3.New(sess)

	// 测试连接
	_, err = client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, wrapStorageError("failed to connect to S3", err)
	}

	return &S3Storage{
		config: config,
		client: client,
	}, nil
}

// Type 返回存储类型
func (s *S3Storage) Type() StorageType {
	return StorageTypeS3
}

// Config 返回存储配置
func (s *S3Storage) Config() StorageConfig {
	return s.config
}

// Save 保存文件到S3
func (s *S3Storage) Save(ctx context.Context, key string, data io.Reader, size int64) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}

	uploader := s3manager.NewUploaderWithClient(s.client)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
		Body:   data,
	})

	if err != nil {
		return wrapStorageError("failed to upload file to S3", err)
	}

	return nil
}

// Get 从S3获取文件
func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// 检查是否是文件不存在错误
		if isNotFoundError(err) {
			return nil, ErrFileNotFound
		}
		return nil, wrapStorageError("failed to get file from S3", err)
	}

	return result.Body, nil
}

// Delete 从S3删除文件
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}

	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return wrapStorageError("failed to delete file from S3", err)
	}

	return nil
}

// Exists 检查文件是否存在于S3
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	if !IsValidKey(key) {
		return false, ErrInvalidKey
	}

	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if isNotFoundError(err) {
			return false, nil
		}
		return false, wrapStorageError("failed to check file existence in S3", err)
	}

	return true, nil
}

// Stat 获取S3文件信息
func (s *S3Storage) Stat(ctx context.Context, key string) (*FileInfo, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	result, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if isNotFoundError(err) {
			return nil, ErrFileNotFound
		}
		return nil, wrapStorageError("failed to get file info from S3", err)
	}

	var lastModified int64
	if result.LastModified != nil {
		lastModified = result.LastModified.Unix()
	}

	return &FileInfo{
		Path:         key,
		Size:         aws.Int64Value(result.ContentLength),
		LastModified: lastModified,
		IsDir:        false, // S3没有目录概念
		MimeType:     aws.StringValue(result.ContentType),
		ETag:         aws.StringValue(result.ETag),
	}, nil
}

// Copy 在S3中复制文件
func (s *S3Storage) Copy(ctx context.Context, srcKey, dstKey string) error {
	if !IsValidKey(srcKey) || !IsValidKey(dstKey) {
		return ErrInvalidKey
	}

	source := fmt.Sprintf("%s/%s", s.config.Bucket, srcKey)
	_, err := s.client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.Bucket),
		CopySource: aws.String(source),
		Key:        aws.String(dstKey),
	})

	if err != nil {
		return wrapStorageError("failed to copy file in S3", err)
	}

	return nil
}

// Move 在S3中移动文件（复制后删除）
func (s *S3Storage) Move(ctx context.Context, srcKey, dstKey string) error {
	if !IsValidKey(srcKey) || !IsValidKey(dstKey) {
		return ErrInvalidKey
	}

	// 先复制
	if err := s.Copy(ctx, srcKey, dstKey); err != nil {
		return err
	}

	// 再删除源文件
	if err := s.Delete(ctx, srcKey); err != nil {
		// 如果删除失败，尝试回滚（删除目标文件）
		s.Delete(ctx, dstKey)
		return err
	}

	return nil
}

// List 列出S3中的文件
func (s *S3Storage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	if prefix != "" && !IsValidKey(prefix) {
		return nil, ErrInvalidKey
	}

	var files []FileInfo
	var continuationToken *string

	for {
		result, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.config.Bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		})

		if err != nil {
			return nil, wrapStorageError("failed to list files in S3", err)
		}

		for _, obj := range result.Contents {
			// 跳过目录标记（S3中没有真正的目录）
			if *obj.Key == prefix || (*obj.Key)[len(*obj.Key)-1:] == "/" {
				continue
			}

			files = append(files, FileInfo{
				Path:         *obj.Key,
				Size:         *obj.Size,
				LastModified: obj.LastModified.Unix(),
				IsDir:        false,
				MimeType:     GetMimeType(*obj.Key),
				ETag:         *obj.ETag,
			})
		}

		// 如果还有更多结果，继续获取
		if result.NextContinuationToken == nil {
			break
		}
		continuationToken = result.NextContinuationToken
	}

	return files, nil
}

// CreateDir 在S3中创建目录（S3没有目录概念，创建空对象作为目录标记）
func (s *S3Storage) CreateDir(ctx context.Context, path string) error {
	if !IsValidKey(path) {
		return ErrInvalidKey
	}

	// 确保路径以斜杠结尾
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	// 创建空对象作为目录标记
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return wrapStorageError("failed to create directory in S3", err)
	}

	return nil
}

// DeleteDir 删除S3中的目录（删除所有以指定前缀开头的对象）
func (s *S3Storage) DeleteDir(ctx context.Context, path string) error {
	if !IsValidKey(path) {
		return ErrInvalidKey
	}

	// 确保路径以斜杠结尾
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	// 列出所有要删除的对象
	var objectsToDelete []*s3.ObjectIdentifier
	var continuationToken *string

	for {
		result, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.config.Bucket),
			Prefix:            aws.String(path),
			ContinuationToken: continuationToken,
		})

		if err != nil {
			return wrapStorageError("failed to list directory contents in S3", err)
		}

		for _, obj := range result.Contents {
			objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		if result.NextContinuationToken == nil {
			break
		}
		continuationToken = result.NextContinuationToken
	}

	// 如果没有对象可删除，直接返回
	if len(objectsToDelete) == 0 {
		return nil
	}

	// 批量删除（S3支持最多1000个对象一次删除）
	for i := 0; i < len(objectsToDelete); i += 1000 {
		end := i + 1000
		if end > len(objectsToDelete) {
			end = len(objectsToDelete)
		}

		batch := objectsToDelete[i:end]
		_, err := s.client.DeleteObjectsWithContext(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.config.Bucket),
			Delete: &s3.Delete{
				Objects: batch,
				Quiet:   aws.Bool(true),
			},
		})

		if err != nil {
			return wrapStorageError("failed to delete directory from S3", err)
		}
	}

	return nil
}

// InitiateMultipartUpload 初始化S3分片上传
func (s *S3Storage) InitiateMultipartUpload(ctx context.Context, key string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	result, err := s.client.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", wrapStorageError("failed to initiate multipart upload in S3", err)
	}

	return *result.UploadId, nil
}

// UploadPart 上传S3分片
func (s *S3Storage) UploadPart(ctx context.Context, uploadID string, partNumber int, data io.Reader) (string, error) {
	// 注意：这里需要知道key，简化实现中我们需要从上下文或其他方式获取
	// 在实际实现中，可能需要存储uploadID和key的映射关系
	return "", fmt.Errorf("not implemented in simplified version")
}

// CompleteMultipartUpload 完成S3分片上传
func (s *S3Storage) CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) error {
	return fmt.Errorf("not implemented in simplified version")
}

// AbortMultipartUpload 中止S3分片上传
func (s *S3Storage) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	// 注意：这里需要知道key，简化实现中我们需要从上下文或其他方式获取
	return fmt.Errorf("not implemented in simplified version")
}

// GetURL 获取S3文件URL
func (s *S3Storage) GetURL(ctx context.Context, key string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(15 * time.Minute) // 15分钟有效期
	if err != nil {
		return "", wrapStorageError("failed to generate S3 URL", err)
	}

	return url, nil
}

// GetDownloadURL 获取S3下载URL
func (s *S3Storage) GetDownloadURL(ctx context.Context, key string, filename string) (string, error) {
	if !IsValidKey(key) {
		return "", ErrInvalidKey
	}

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(s.config.Bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	})

	url, err := req.Presign(15 * time.Minute) // 15分钟有效期
	if err != nil {
		return "", wrapStorageError("failed to generate S3 download URL", err)
	}

	return url, nil
}

// 辅助函数

// isNotFoundError 检查是否是文件不存在错误
func isNotFoundError(err error) bool {
	// 这里简化实现，实际应该检查具体的AWS错误类型
	return err != nil && (err.Error() == "NoSuchKey" || err.Error() == "NotFound")
}

// MinIOStorage MinIO存储实现（继承S3Storage）
type MinIOStorage struct {
	*S3Storage
}

// NewMinIOStorage 创建MinIO存储实例
func NewMinIOStorage(config StorageConfig) (*MinIOStorage, error) {
	// MinIO使用S3兼容API，所以重用S3Storage
	s3Storage, err := NewS3Storage(config)
	if err != nil {
		return nil, err
	}

	return &MinIOStorage{s3Storage}, nil
}

// Type 返回存储类型
func (m *MinIOStorage) Type() StorageType {
	return StorageTypeMinIO
}