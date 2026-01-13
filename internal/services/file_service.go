package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloud-storage/internal/config"
	"cloud-storage/internal/models"
	"cloud-storage/internal/pkg/storage"
	"cloud-storage/internal/repositories"
)

// FileService 文件服务
type FileService struct {
	cfg             *config.Config
	db              *gorm.DB
	fileRepo        repositories.FileRepository
	userRepo        repositories.UserRepository
	fileVersionRepo repositories.FileVersionRepository
	storage         storage.Storage
}

// NewFileService 创建文件服务实例
func NewFileService(
	cfg *config.Config,
	db *gorm.DB,
	fileRepo repositories.FileRepository,
	userRepo repositories.UserRepository,
	storage storage.Storage,
) *FileService {
	return &FileService{
		cfg:             cfg,
		db:              db,
		fileRepo:        fileRepo,
		userRepo:        userRepo,
		fileVersionRepo: repositories.NewFileVersionRepository(db),
		storage:         storage,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(
	ctx *gin.Context,
	userID uuid.UUID,
	fileHeader *multipart.FileHeader,
	req models.FileUploadRequest,
) (*models.File, error) {
	// 检查用户存储配额
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 检查配额
	if !user.CheckStorageQuota(fileHeader.Size) {
		return nil, fmt.Errorf("storage quota exceeded")
	}

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// 生成文件信息
	filename := fileHeader.Filename
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = storage.GetMimeType(filename)
	}

	// 检查文件是否已存在
	existingFile, err := s.fileRepo.FindByUserAndName(userID, req.ParentID, filename)
	if err == nil && existingFile != nil {
		if req.Override {
			// 覆盖现有文件
			return s.updateExistingFile(ctx, userID, existingFile, file, fileHeader.Size, mimeType)
		}
		return nil, fmt.Errorf("file already exists")
	}

	// 创建文件记录
	newFile := &models.File{
		UserID:   userID,
		ParentID: req.ParentID,
		Name:     filename,
		Size:     fileHeader.Size,
		MimeType: mimeType,
		Type:     models.FileTypeFile,
		IsPublic: req.IsPublic,
		Version:  1,
	}

	// 在事务中保存文件
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 保存文件记录
	if err := s.fileRepo.CreateWithTx(tx, newFile); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	// 保存文件内容到存储
	storageKey := storage.GenerateFileKey(userID, newFile.Path)
	if err := s.storage.Save(ctx, storageKey, file, fileHeader.Size); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save file to storage: %w", err)
	}

	// 更新用户已使用存储
	if err := user.UpdateUsedStorage(tx, fileHeader.Size); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update user storage: %w", err)
	}

	// 创建文件版本记录
	fileVersion := &models.FileVersion{
		FileID:        newFile.ID,
		VersionNumber: 1,
		FileSize:      fileHeader.Size,
		FileHash:      "", // 可以计算文件哈希
		StoragePath:   storageKey,
		MimeType:      mimeType,
		CreatedBy:     userID,
	}

	if err := tx.Create(fileVersion).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create file version: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newFile, nil
}

// updateExistingFile 更新现有文件
func (s *FileService) updateExistingFile(
	ctx *gin.Context,
	userID uuid.UUID,
	existingFile *models.File,
	file io.Reader,
	size int64,
	mimeType string,
) (*models.File, error) {
	// 计算存储空间变化
	sizeDelta := size - existingFile.Size

	// 检查用户存储配额
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.CheckStorageQuota(sizeDelta) {
		return nil, fmt.Errorf("storage quota exceeded")
	}

	// 在事务中更新文件
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 更新文件记录
	existingFile.Size = size
	existingFile.MimeType = mimeType
	existingFile.Version++

	updates := map[string]interface{}{
		"size":      size,
		"mime_type": mimeType,
		"version":   existingFile.Version,
	}

	if err := s.fileRepo.UpdateWithTx(tx, existingFile.ID, updates); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update file record: %w", err)
	}

	// 保存新版本到存储
	storageKey := storage.GenerateFileKey(userID, existingFile.Path)
	if err := s.storage.Save(ctx, storageKey, file, size); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save file to storage: %w", err)
	}

	// 更新用户已使用存储
	if err := user.UpdateUsedStorage(tx, sizeDelta); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update user storage: %w", err)
	}

	// 创建新版本记录
	fileVersion := &models.FileVersion{
		FileID:        existingFile.ID,
		VersionNumber: existingFile.Version,
		FileSize:      size,
		FileHash:      "", // 可以计算文件哈希
		StoragePath:   storageKey,
		MimeType:      mimeType,
		CreatedBy:     userID,
	}

	if err := tx.Create(fileVersion).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create file version: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return existingFile, nil
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(
	ctx *gin.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
) (io.ReadCloser, *models.File, error) {
	// 获取文件信息
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID && !file.IsPublic {
		return nil, nil, fmt.Errorf("permission denied")
	}

	// 获取文件内容
	storageKey := storage.GenerateFileKey(file.UserID, file.Path)
	reader, err := s.storage.Get(ctx, storageKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file from storage: %w", err)
	}

	return reader, file, nil
}

// CreateDirectory 创建目录
func (s *FileService) CreateDirectory(
	ctx *gin.Context,
	userID uuid.UUID,
	req models.FileCreateRequest,
) (*models.File, error) {
	// 验证请求
	if req.Type != models.FileTypeDir {
		return nil, fmt.Errorf("invalid file type for directory creation")
	}

	// 检查目录是否已存在
	existingDir, err := s.fileRepo.FindByUserAndName(userID, req.ParentID, req.Name)
	if err == nil && existingDir != nil {
		return nil, fmt.Errorf("directory already exists")
	}

	// 创建目录记录
	directory := &models.File{
		UserID:   userID,
		ParentID: req.ParentID,
		Name:     req.Name,
		Size:     0,
		Type:     models.FileTypeDir,
		IsPublic: req.IsPublic,
		Version:  1,
	}

	// 保存目录记录
	if err := s.fileRepo.Create(directory); err != nil {
		return nil, fmt.Errorf("failed to create directory record: %w", err)
	}

	// 在存储中创建目录
	storageKey := storage.GenerateFileKey(userID, directory.Path)
	if err := s.storage.CreateDir(ctx, storageKey); err != nil {
		// 如果存储创建失败，删除数据库记录
		s.fileRepo.Delete(directory.ID)
		return nil, fmt.Errorf("failed to create directory in storage: %w", err)
	}

	return directory, nil
}

// GetFileList 获取文件列表
func (s *FileService) GetFileList(
	userID uuid.UUID,
	filter models.FileFilter,
) ([]models.File, int64, error) {
	// 设置用户ID过滤器
	filter.UserID = &userID

	// 获取文件列表
	files, err := s.fileRepo.FindAll(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get file list: %w", err)
	}

	// 获取总数
	total, err := s.fileRepo.Count(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
	}

	return files, total, nil
}

// GetFileByID 根据ID获取文件
func (s *FileService) GetFileByID(
	userID uuid.UUID,
	fileID uuid.UUID,
) (*models.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID && !file.IsPublic {
		return nil, fmt.Errorf("permission denied")
	}

	return file, nil
}

// UpdateFile 更新文件信息
func (s *FileService) UpdateFile(
	userID uuid.UUID,
	fileID uuid.UUID,
	req models.FileUpdateRequest,
) (*models.File, error) {
	// 获取文件
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// 更新文件信息
	updates := make(map[string]interface{})

	if req.Name != nil {
		// 检查新名称是否已存在
		existingFile, err := s.fileRepo.FindByUserAndName(userID, file.ParentID, *req.Name)
		if err == nil && existingFile != nil && existingFile.ID != fileID {
			return nil, fmt.Errorf("file with this name already exists")
		}
		updates["name"] = *req.Name
	}

	if req.ParentID != nil {
		// 检查目标目录是否存在且不是当前文件的子目录
		if *req.ParentID != file.ID {
			targetDir, err := s.fileRepo.FindByID(*req.ParentID)
			if err != nil || targetDir.Type != models.FileTypeDir {
				return nil, fmt.Errorf("invalid target directory")
			}

			// 检查是否移动到自己的子目录
			if s.isDescendant(file.ID, *req.ParentID) {
				return nil, fmt.Errorf("cannot move directory into its own subdirectory")
			}
		}
		updates["parent_id"] = *req.ParentID
	}

	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	// 应用更新
	if err := s.fileRepo.Update(fileID, updates); err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	// 重新加载文件信息
	updatedFile, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload file: %w", err)
	}

	return updatedFile, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(
	ctx *gin.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
	permanent bool,
) error {
	// 获取文件
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return fmt.Errorf("permission denied")
	}

	if permanent {
		// 永久删除
		return s.permanentDeleteFile(ctx, userID, file)
	}

	// 软删除
	return s.softDeleteFile(file)
}

// permanentDeleteFile 永久删除文件
func (s *FileService) permanentDeleteFile(
	ctx *gin.Context,
	userID uuid.UUID,
	file *models.File,
) error {
	// 在事务中删除文件
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if file.Type == models.FileTypeDir {
		// 递归删除目录下的所有文件
		if err := s.deleteDirectoryRecursive(ctx, tx, userID, file); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete directory: %w", err)
		}
	} else {
		// 删除单个文件
		if err := s.deleteSingleFile(ctx, tx, userID, file); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// deleteDirectoryRecursive 递归删除目录
func (s *FileService) deleteDirectoryRecursive(
	ctx *gin.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	directory *models.File,
) error {
	// 获取目录下的所有文件
	filter := models.FileFilter{
		UserID:   &userID,
		ParentID: &directory.ID,
		Deleted:  &[]bool{false}[0],
	}

	files, err := s.fileRepo.FindAllWithTx(tx, filter)
	if err != nil {
		return err
	}

	// 递归删除子文件和子目录
	for _, file := range files {
		if file.Type == models.FileTypeDir {
			if err := s.deleteDirectoryRecursive(ctx, tx, userID, &file); err != nil {
				return err
			}
		} else {
			if err := s.deleteSingleFile(ctx, tx, userID, &file); err != nil {
				return err
			}
		}
	}

	// 删除目录记录
	if err := s.fileRepo.DeleteWithTx(tx, directory.ID); err != nil {
		return err
	}

	// 删除存储中的目录
	storageKey := storage.GenerateFileKey(userID, directory.Path)
	if err := s.storage.DeleteDir(ctx, storageKey); err != nil {
		return err
	}

	return nil
}

// deleteSingleFile 删除单个文件
func (s *FileService) deleteSingleFile(
	ctx *gin.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	file *models.File,
) error {
	// 删除文件记录
	if err := s.fileRepo.DeleteWithTx(tx, file.ID); err != nil {
		return err
	}

	// 删除存储中的文件
	storageKey := storage.GenerateFileKey(userID, file.Path)
	if err := s.storage.Delete(ctx, storageKey); err != nil {
		return err
	}

	// 更新用户已使用存储
	user, err := s.userRepo.FindByIDWithTx(tx, userID)
	if err != nil {
		return err
	}

	if err := user.UpdateUsedStorage(tx, -file.Size); err != nil {
		return err
	}

	return nil
}

// softDeleteFile 软删除文件
func (s *FileService) softDeleteFile(file *models.File) error {
	// 软删除文件记录
	return s.fileRepo.SoftDelete(file.ID)
}

// MoveFile 移动文件
func (s *FileService) MoveFile(
	ctx *gin.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
	req models.FileMoveRequest,
) (*models.File, error) {
	// 获取文件
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// 检查目标目录
	targetDir, err := s.fileRepo.FindByID(*req.TargetParentID)
	if err != nil || targetDir.Type != models.FileTypeDir {
		return nil, fmt.Errorf("invalid target directory")
	}

	// 检查是否移动到自己的子目录
	if file.Type == models.FileTypeDir && s.isDescendant(file.ID, *req.TargetParentID) {
		return nil, fmt.Errorf("cannot move directory into its own subdirectory")
	}

	// 检查目标位置是否已存在同名文件
	existingFile, err := s.fileRepo.FindByUserAndName(userID, req.TargetParentID, file.Name)
	if err == nil && existingFile != nil {
		return nil, fmt.Errorf("file with this name already exists in target directory")
	}

	// 在事务中移动文件
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 更新文件父目录
	updates := map[string]interface{}{
		"parent_id": req.TargetParentID,
	}

	if err := s.fileRepo.UpdateWithTx(tx, fileID, updates); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	// 如果文件是目录，需要更新所有子文件的路径
	if file.Type == models.FileTypeDir {
		if err := s.updateDescendantPaths(tx, file); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update descendant paths: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新加载文件信息
	updatedFile, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload file: %w", err)
	}

	return updatedFile, nil
}

// CopyFile 复制文件
func (s *FileService) CopyFile(
	ctx *gin.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
	req models.FileCopyRequest,
) (*models.File, error) {
	// 获取源文件
	sourceFile, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if sourceFile.UserID != userID && !sourceFile.IsPublic {
		return nil, fmt.Errorf("permission denied")
	}

	// 检查目标目录
	targetDir, err := s.fileRepo.FindByID(*req.TargetParentID)
	if err != nil || targetDir.Type != models.FileTypeDir {
		return nil, fmt.Errorf("invalid target directory")
	}

	// 确定新文件名
	newName := sourceFile.Name
	if req.NewName != nil {
		newName = *req.NewName
	}

	// 检查目标位置是否已存在同名文件
	existingFile, err := s.fileRepo.FindByUserAndName(userID, req.TargetParentID, newName)
	if err == nil && existingFile != nil {
		return nil, fmt.Errorf("file with this name already exists in target directory")
	}

	// 检查用户存储配额
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.CheckStorageQuota(sourceFile.Size) {
		return nil, fmt.Errorf("storage quota exceeded")
	}

	// 在事务中复制文件
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 创建文件副本
	copiedFile, err := s.copyFileRecursive(ctx, tx, userID, sourceFile, req.TargetParentID, newName)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 更新用户已使用存储
	if err := user.UpdateUsedStorage(tx, sourceFile.Size); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update user storage: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return copiedFile, nil
}

// copyFileRecursive 递归复制文件
func (s *FileService) copyFileRecursive(
	ctx *gin.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	sourceFile *models.File,
	targetParentID *uuid.UUID,
	newName string,
) (*models.File, error) {
	// 创建文件记录副本
	copiedFile := &models.File{
		UserID:   userID,
		ParentID: targetParentID,
		Name:     newName,
		Size:     sourceFile.Size,
		MimeType: sourceFile.MimeType,
		Type:     sourceFile.Type,
		IsPublic: sourceFile.IsPublic,
		Version:  1,
	}

	// 保存文件记录
	if err := s.fileRepo.CreateWithTx(tx, copiedFile); err != nil {
		return nil, err
	}

	if sourceFile.Type == models.FileTypeFile {
		// 复制文件内容
		srcStorageKey := storage.GenerateFileKey(sourceFile.UserID, sourceFile.Path)
		dstStorageKey := storage.GenerateFileKey(userID, copiedFile.Path)

		// 获取源文件
		reader, err := s.storage.Get(ctx, srcStorageKey)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		// 保存副本
		if err := s.storage.Save(ctx, dstStorageKey, reader, sourceFile.Size); err != nil {
			return nil, err
		}

		// 创建版本记录
		fileVersion := &models.FileVersion{
			FileID:        copiedFile.ID,
			VersionNumber: 1,
			FileSize:      sourceFile.Size,
			FileHash:      "", // 可以计算文件哈希
			StoragePath:   dstStorageKey,
			MimeType:      sourceFile.MimeType,
			CreatedBy:     userID,
		}

		if err := tx.Create(fileVersion).Error; err != nil {
			return nil, err
		}
	} else if sourceFile.Type == models.FileTypeDir {
		// 在存储中创建目录
		dstStorageKey := storage.GenerateFileKey(userID, copiedFile.Path)
		if err := s.storage.CreateDir(ctx, dstStorageKey); err != nil {
			return nil, err
		}

		// 递归复制子文件
		filter := models.FileFilter{
			UserID:   &sourceFile.UserID,
			ParentID: &sourceFile.ID,
			Deleted:  &[]bool{false}[0],
		}

		childFiles, err := s.fileRepo.FindAllWithTx(tx, filter)
		if err != nil {
			return nil, err
		}

		for _, childFile := range childFiles {
			_, err := s.copyFileRecursive(ctx, tx, userID, &childFile, &copiedFile.ID, childFile.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return copiedFile, nil
}

// GetFileVersions 获取文件版本列表
func (s *FileService) GetFileVersions(
	userID uuid.UUID,
	fileID uuid.UUID,
) ([]models.FileVersion, error) {
	// 获取文件
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// 获取版本列表
	versions, err := s.fileVersionRepo.FindByFileID(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file versions: %w", err)
	}

	return versions, nil
}

// RestoreFileVersion 恢复文件版本
func (s *FileService) RestoreFileVersion(
	ctx *gin.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
	versionNumber int,
) (*models.File, error) {
	// 获取文件
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// 获取指定版本
	version, err := s.fileVersionRepo.FindByVersion(fileID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// 创建新版本记录
	newVersion := &models.FileVersion{
		FileID:        fileID,
		VersionNumber: file.Version + 1,
		FileSize:      file.Size,
		FileHash:      file.Hash,
		StoragePath:   file.Path,
		MimeType:      file.MimeType,
		CreatedBy:     userID,
	}

	// 在事务中恢复
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 保存当前版本
	if err := s.fileVersionRepo.Create(newVersion); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save current version: %w", err)
	}

	// 从指定版本恢复文件内容
	srcStorageKey := version.StoragePath
	dstStorageKey := storage.GenerateFileKey(userID, file.Path)

	reader, err := s.storage.Get(ctx, srcStorageKey)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get version file: %w", err)
	}
	defer reader.Close()

	if err := s.storage.Save(ctx, dstStorageKey, reader, version.FileSize); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to restore file: %w", err)
	}

	// 更新文件信息
	updates := map[string]interface{}{
		"size":         version.FileSize,
		"mime_type":    version.MimeType,
		"hash":         version.FileHash,
		"version":      file.Version + 1,
		"storage_path": version.StoragePath,
	}

	if err := s.fileRepo.UpdateWithTx(tx, fileID, updates); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新加载文件信息
	return s.fileRepo.FindByID(fileID)
}

// SearchFiles 搜索文件
func (s *FileService) SearchFiles(
	userID uuid.UUID,
	query string,
	searchIn string,
	page, pageSize int,
) ([]models.File, int64, error) {
	// 构建搜索条件
	filter := models.FileFilter{
		UserID:   &userID,
		Deleted:  &[]bool{false}[0],
		Page:     page,
		PageSize: pageSize,
	}

	// 根据搜索类型设置不同的条件
	switch searchIn {
	case "name":
		filter.Name = &query
	case "path":
		// 路径搜索需要特殊处理
		// 这里简化实现
		filter.Name = &query
	case "content":
		// 内容搜索需要全文索引，这里简化实现
		filter.Name = &query
	}

	// 搜索文件
	files, err := s.fileRepo.FindAll(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search files: %w", err)
	}

	// 获取总数
	total, err := s.fileRepo.Count(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return files, total, nil
}

// GetFileStats 获取文件统计信息
func (s *FileService) GetFileStats(userID uuid.UUID) (*models.FileStats, error) {
	stats := &models.FileStats{}

	// 获取用户文件统计
	// 这里需要实现具体的统计查询
	// 暂时返回空数据

	return stats, nil
}

// GetRecycledFiles 获取回收站文件
func (s *FileService) GetRecycledFiles(
	userID uuid.UUID,
	page, pageSize int,
) ([]models.File, int64, error) {
	filter := models.FileFilter{
		UserID:   &userID,
		Deleted:  &[]bool{true}[0],
		Page:     page,
		PageSize: pageSize,
	}

	files, err := s.fileRepo.FindAll(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get recycled files: %w", err)
	}

	total, err := s.fileRepo.Count(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count recycled files: %w", err)
	}

	return files, total, nil
}

// RestoreRecycledFile 恢复回收站文件
func (s *FileService) RestoreRecycledFile(
	userID uuid.UUID,
	fileID uuid.UUID,
) error {
	// 获取文件（包括已删除的）
	file, err := s.fileRepo.FindByIDIncludingDeleted(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// 检查权限
	if file.UserID != userID {
		return fmt.Errorf("permission denied")
	}

	// 恢复文件
	return s.fileRepo.Restore(fileID)
}

// CleanupRecycledFiles 清理回收站文件
func (s *FileService) CleanupRecycledFiles(
	ctx *gin.Context,
	userID uuid.UUID,
	daysOld int,
) (int, error) {
	// 计算截止日期
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)

	// 获取需要清理的文件
	files, err := s.fileRepo.FindOldRecycledFiles(userID, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to find old recycled files: %w", err)
	}

	// 永久删除文件
	deletedCount := 0
	for _, file := range files {
		if err := s.permanentDeleteFile(ctx, userID, &file); err != nil {
			// 记录错误但继续处理其他文件
			continue
		}
		deletedCount++
	}

	return deletedCount, nil
}

// 辅助方法

// isDescendant 检查一个文件是否是另一个文件的后代
func (s *FileService) isDescendant(fileID, potentialAncestorID uuid.UUID) bool {
	if fileID == potentialAncestorID {
		return true
	}

	// 获取文件的父目录
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return false
	}

	if file.ParentID == nil {
		return false
	}

	// 递归检查父目录
	return s.isDescendant(*file.ParentID, potentialAncestorID)
}

// updateDescendantPaths 更新后代文件的路径
func (s *FileService) updateDescendantPaths(tx *gorm.DB, directory *models.File) error {
	// 获取所有子文件
	filter := models.FileFilter{
		UserID:   &directory.UserID,
		ParentID: &directory.ID,
		Deleted:  &[]bool{false}[0],
	}

	childFiles, err := s.fileRepo.FindAllWithTx(tx, filter)
	if err != nil {
		return err
	}

	// 递归更新子文件路径
	for _, childFile := range childFiles {
		// 更新路径（GORM的BeforeUpdate钩子会自动处理）
		if err := tx.Save(&childFile).Error; err != nil {
			return err
		}

		// 如果子文件是目录，递归更新
		if childFile.Type == models.FileTypeDir {
			if err := s.updateDescendantPaths(tx, &childFile); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetStorageUsage 获取存储使用情况
func (s *FileService) GetStorageUsage(userID uuid.UUID) (int64, int64, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get user: %w", err)
	}

	return user.UsedStorage, user.StorageQuota, nil
}

// GenerateShareToken 生成分享令牌
func (s *FileService) GenerateShareToken(fileID uuid.UUID) (string, error) {
	// 生成随机令牌
	token := uuid.New().String()
	token = strings.ReplaceAll(token, "-", "")[:32]

	// 检查令牌是否唯一（概率极低，但需要检查）
	existingFile, err := s.fileRepo.FindByShareToken(token)
	if err == nil && existingFile != nil {
		// 如果冲突，重新生成
		return s.GenerateShareToken(fileID)
	}

	return token, nil
}
