package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloud-storage/internal/models"
	"cloud-storage/internal/services"
)

type ShareHandler struct {
	shareService *services.ShareService
}

func NewShareHandler(shareService *services.ShareService) *ShareHandler {
	return &ShareHandler{
		shareService: shareService,
	}
}

func (h *ShareHandler) RegisterRoutes(protected *gin.RouterGroup, public *gin.RouterGroup) {
	shares := protected.Group("/shares")
	{
		shares.POST("", h.CreateShare)
		shares.GET("", h.GetUserShares)
		shares.GET("/:id", h.GetShare)
		shares.PUT("/:id", h.UpdateShare)
		shares.DELETE("/:id", h.DeleteShare)
		shares.POST("/batch-delete", h.BatchDeleteShares)
		shares.GET("/stats", h.GetShareStats)
	}

	publicRoutes := public.Group("/s")
	{
		publicRoutes.GET("/:token", h.AccessShare)
		publicRoutes.GET("/:token/download", h.DownloadSharedFile)
	}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.ShareCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share, err := h.shareService.CreateShare(userID, req.FileID, req)
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

	response := share.ToResponse()
	response.ShareURL = getShareURL(c, share.ShareToken)

	c.JSON(http.StatusCreated, response)
}

func (h *ShareHandler) GetUserShares(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var filter models.ShareFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.UserIDStr != "" {
		userID, err := uuid.Parse(filter.UserIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		filter.UserID = &userID
	}

	if filter.FileIDStr != "" {
		fileID, err := uuid.Parse(filter.FileIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id format"})
			return
		}
		filter.FileID = &fileID
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	shares, total, err := h.shareService.GetUserShares(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.ShareResponse
	for _, share := range shares {
		r := share.ToResponse()
		r.ShareURL = getShareURL(c, share.ShareToken)
		response = append(response, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"shares": response,
		"total":  total,
		"page":   filter.Page,
		"size":   filter.PageSize,
	})
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
		return
	}

	share, err := h.shareService.GetShare(shareID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "share not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := share.ToResponse()
	response.ShareURL = getShareURL(c, share.ShareToken)

	c.JSON(http.StatusOK, response)
}

func (h *ShareHandler) UpdateShare(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
		return
	}

	var req models.ShareUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share, err := h.shareService.UpdateShare(shareID, userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "share not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := share.ToResponse()
	response.ShareURL = getShareURL(c, share.ShareToken)

	c.JSON(http.StatusOK, response)
}

func (h *ShareHandler) DeleteShare(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	shareID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
		return
	}

	err = h.shareService.DeleteShare(shareID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "share not found" {
			status = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "share deleted successfully"})
}

func (h *ShareHandler) BatchDeleteShares(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.ShareBulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deletedCount, err := h.shareService.BatchDeleteShares(req.ShareIDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "shares deleted successfully",
		"deleted_count": deletedCount,
	})
}

func (h *ShareHandler) GetShareStats(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	stats, err := h.shareService.GetShareStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ShareHandler) AccessShare(c *gin.Context) {
	token := c.Param("token")

	var password *string
	if c.Query("password") != "" {
		pw := c.Query("password")
		password = &pw
	}

	share, err := h.shareService.AccessShare(token, password)
	if err != nil {
		status := http.StatusForbidden
		if err.Error() == "share not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := share.ToResponse()
	response.ShareURL = getShareURL(c, share.ShareToken)

	if share.FileID != uuid.Nil {
		file := share.File.ToResponse()
		response.FileName = file.Name
		response.FileSize = file.Size
		response.FileType = string(file.Type)
	}

	c.JSON(http.StatusOK, gin.H{
		"share": response,
	})
}

func (h *ShareHandler) DownloadSharedFile(c *gin.Context) {
	token := c.Param("token")

	var password *string
	if c.Query("password") != "" {
		pw := c.Query("password")
		password = &pw
	}

	file, err := h.shareService.DownloadSharedFile(token, password)
	if err != nil {
		status := http.StatusForbidden
		if err.Error() == "share not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file":         file.ToResponse(),
		"download_url": c.Request.Host + "/api/v1/s/" + token + "/download",
	})
}

func getShareURL(c *gin.Context, token string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + "/api/v1/s/" + token
}
