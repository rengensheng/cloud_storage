package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloud-storage/internal/models"
	"cloud-storage/internal/repositories"
	"cloud-storage/internal/services"
)

type OperationLogHandler struct {
	logService *services.OperationLogService
}

func NewOperationLogHandler(logService *services.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{
		logService: logService,
	}
}

func (h *OperationLogHandler) RegisterRoutes(router *gin.RouterGroup) {
	logs := router.Group("/logs")
	{
		logs.GET("", h.GetLogs)
		logs.GET("/stats", h.GetLogStats)
		logs.DELETE("/cleanup", h.CleanupLogs)
	}
}

func (h *OperationLogHandler) GetLogs(c *gin.Context) {
	var filter models.OperationLogFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 50
	}

	logs, total, err := h.logService.GetLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.OperationLogResponse
	for _, log := range logs {
		response = append(response, log.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  response,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

func (h *OperationLogHandler) GetLogStats(c *gin.Context) {
	userIDStr := c.Query("user_id")
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var startDate, endDate time.Time
		if startDateStr != "" {
			startDate, err = time.Parse(time.RFC3339, startDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
				return
			}
		} else {
			startDate = time.Now().AddDate(0, 0, -7)
		}

		if endDateStr != "" {
			endDate, err = time.Parse(time.RFC3339, endDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
				return
			}
		} else {
			endDate = time.Now()
		}

		stats, err := h.logService.GetUserOperationStats(userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":    userID,
			"start_date": startDate,
			"end_date":   endDate,
			"stats":      stats,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
}

func (h *OperationLogHandler) CleanupLogs(c *gin.Context) {
	userRole := c.GetString("role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can cleanup logs"})
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 1 {
		days = 90
	}

	deletedCount, err := h.logService.CleanupOldLogs(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "old logs cleaned up successfully",
		"deleted_count":   deletedCount,
		"days_older_than": days,
	})
}

type AdminHandler struct {
	userRepo     repositories.UserRepository
	logService   *services.OperationLogService
	shareService *services.ShareService
	fileService  *services.FileService
}

func NewAdminHandler(
	userRepo repositories.UserRepository,
	logService *services.OperationLogService,
	shareService *services.ShareService,
	fileService *services.FileService,
) *AdminHandler {
	return &AdminHandler{
		userRepo:     userRepo,
		logService:   logService,
		shareService: shareService,
		fileService:  fileService,
	}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	{
		admin.GET("/stats", h.GetSystemStats)
		admin.GET("/users", h.ListUsers)
		admin.GET("/users/:id", h.GetUser)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
		admin.POST("/users/:id/activate", h.ActivateUser)
		admin.POST("/users/:id/deactivate", h.DeactivateUser)
	}
}

func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.logService.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := models.UserFilter{
		Page:     page,
		PageSize: pageSize,
	}

	users, err := h.userRepo.FindAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	total, _ := h.userRepo.Count(filter)

	c.JSON(http.StatusOK, gin.H{
		"users": response,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.StorageQuota != nil {
		updates["storage_quota"] = *req.StorageQuota
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := h.userRepo.Update(userID, updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    user.ToResponse(),
	})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userRepo.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *AdminHandler) ActivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	updates := map[string]interface{}{"is_active": true}
	if err := h.userRepo.Update(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user activated successfully"})
}

func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	updates := map[string]interface{}{"is_active": false}
	if err := h.userRepo.Update(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated successfully"})
}
