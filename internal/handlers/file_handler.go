package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloud-storage/internal/models"
	"cloud-storage/internal/services"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *services.FileService
}

// NewFileHandler 创建文件处理器实例
func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// RegisterRoutes 注册文件路由
func (h *FileHandler) RegisterRoutes(router *gin.RouterGroup) {
	files := router.Group("/files")
	{
		files.GET("", h.GetFileList)
		files.POST("", h.CreateFileOrDirectory)
		files.GET("/:id", h.GetFile)
		files.PUT("/:id", h.UpdateFile)
		files.DELETE("/:id", h.DeleteFile)
		files.POST("/:id/copy", h.CopyFile)
		files.POST("/:id/move", h.MoveFile)
		files.GET("/:id/download", h.DownloadFile)
		files.GET("/:id/versions", h.GetFileVersions)
		files.POST("/:id/restore-version", h.RestoreFileVersion)
	}

	upload := router.Group("/upload")
	{
		upload.POST("", h.UploadFile)
		upload.POST("/chunk", h.UploadChunk)
	}

	recycle := router.Group("/recycle")
	{
		recycle.GET("", h.GetRecycledFiles)
		recycle.POST("/:id/restore", h.RestoreRecycledFile)
		recycle.DELETE("/cleanup", h.CleanupRecycledFiles)
	}

	search := router.Group("/search")
	{
		search.GET("", h.SearchFiles)
	}

	stats := router.Group("/stats")
	{
		stats.GET("/storage", h.GetStorageUsage)
		stats.GET("/files", h.GetFileStats)
	}
}

// GetFileList 获取文件列表
func (h *FileHandler) GetFileList(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var filter models.FileFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.ParentIDStr != "" {
		parentID, err := uuid.Parse(filter.ParentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id format"})
			return
		}
		filter.ParentID = &parentID
	}

	// 设置默认值
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	files, total, err := h.fileService.GetFileList(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var response []models.FileResponse
	for _, file := range files {
		response = append(response, file.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// CreateFileOrDirectory 创建文件或目录
func (h *FileHandler) CreateFileOrDirectory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.FileCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result *models.File
	var err error

	if req.Type == models.FileTypeDir {
		// 创建目录
		result, err = h.fileService.CreateDirectory(c, userID, req)
	} else {
		// 创建文件需要上传，这里只处理元数据创建
		c.JSON(http.StatusBadRequest, gin.H{"error": "use upload endpoint for file creation"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result.ToResponse())
}

// GetFile 获取文件信息
func (h *FileHandler) GetFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	file, err := h.fileService.GetFileByID(userID, fileID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file.ToResponse())
}

// UpdateFile 更新文件信息
func (h *FileHandler) UpdateFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	var req models.FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.UpdateFile(userID, fileID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		} else if err.Error() == "file with this name already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file.ToResponse())
}

// DeleteFile 删除文件
func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	// 检查是否永久删除
	permanent := c.Query("permanent") == "true"

	err = h.fileService.DeleteFile(c, userID, fileID, permanent)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if permanent {
		c.JSON(http.StatusOK, gin.H{"message": "file permanently deleted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "file moved to recycle bin"})
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// 解析表单数据
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	var req models.FileUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ParentIDStr != "" {
		parentID, err := uuid.Parse(req.ParentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id format"})
			return
		}
		req.ParentID = &parentID
	}

	file, err := h.fileService.UploadFile(c, userID, fileHeader, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "storage quota exceeded" {
			status = http.StatusForbidden
		} else if err.Error() == "file already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, file.ToResponse())
}

// DownloadFile 下载文件
func (h *FileHandler) DownloadFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	reader, file, err := h.fileService.DownloadFile(c, userID, fileID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// 流式传输文件
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, reader)
		return err == nil
	})
}

// CopyFile 复制文件
func (h *FileHandler) CopyFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	var req models.FileCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.CopyFile(c, userID, fileID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		} else if err.Error() == "storage quota exceeded" {
			status = http.StatusForbidden
		} else if err.Error() == "file with this name already exists in target directory" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file.ToResponse())
}

// MoveFile 移动文件
func (h *FileHandler) MoveFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	var req models.FileMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.MoveFile(c, userID, fileID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		} else if err.Error() == "file with this name already exists in target directory" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file.ToResponse())
}

// GetFileVersions 获取文件版本列表
func (h *FileHandler) GetFileVersions(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	versions, err := h.fileService.GetFileVersions(userID, fileID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var response []models.FileVersionResponse
	for _, version := range versions {
		response = append(response, version.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"versions": response})
}

// RestoreFileVersion 恢复文件版本
func (h *FileHandler) RestoreFileVersion(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	var req models.VersionRestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.fileService.RestoreFileVersion(c, userID, fileID, req.VersionNumber)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		} else if err.Error() == "version not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file.ToResponse())
}

// UploadChunk 分片上传
func (h *FileHandler) UploadChunk(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "chunk upload functionality requires additional implementation"})
}

// GetRecycledFiles 获取回收站文件
func (h *FileHandler) GetRecycledFiles(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	files, total, err := h.fileService.GetRecycledFiles(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var response []models.FileResponse
	for _, file := range files {
		response = append(response, file.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// RestoreRecycledFile 恢复回收站文件
func (h *FileHandler) RestoreRecycledFile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	err = h.fileService.RestoreRecycledFile(userID, fileID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "file not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file restored successfully"})
}

// CleanupRecycledFiles 清理回收站文件
func (h *FileHandler) CleanupRecycledFiles(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	deletedCount, err := h.fileService.CleanupRecycledFiles(c, userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "recycled files cleaned up",
		"deleted_count": deletedCount,
	})
}

// SearchFiles 搜索文件
func (h *FileHandler) SearchFiles(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	searchIn := c.DefaultQuery("search_in", "name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	files, total, err := h.fileService.SearchFiles(userID, query, searchIn, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应格式
	var response []models.FileResponse
	for _, file := range files {
		response = append(response, file.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"total": total,
		"page":  page,
		"size":  pageSize,
		"query": query,
	})
}

// GetStorageUsage 获取存储使用情况
func (h *FileHandler) GetStorageUsage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	used, quota, err := h.fileService.GetStorageUsage(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usagePercent := 0.0
	if quota > 0 {
		usagePercent = float64(used) / float64(quota) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"used":          used,
		"quota":         quota,
		"available":     quota - used,
		"usage_percent": usagePercent,
		"usage_readable": fmt.Sprintf("%s / %s",
			formatFileSize(used),
			formatFileSize(quota)),
	})
}

// GetFileStats 获取文件统计信息
func (h *FileHandler) GetFileStats(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	stats, err := h.fileService.GetFileStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// 辅助函数

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// ShareFile 分享文件（需要分享服务）
func (h *FileHandler) ShareFile(c *gin.Context) {
	// 分享功能需要分享服务
	c.JSON(http.StatusNotImplemented, gin.H{"error": "share functionality not implemented yet"})
}

// GetSharedFile 获取分享的文件
func (h *FileHandler) GetSharedFile(c *gin.Context) {
	// 分享功能需要分享服务
	c.JSON(http.StatusNotImplemented, gin.H{"error": "share functionality not implemented yet"})
}
