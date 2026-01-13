package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloud-storage/internal/models"
	"cloud-storage/internal/repositories"
)

type OperationLogService struct {
	logRepo repositories.OperationLogRepository
}

func NewOperationLogService(logRepo repositories.OperationLogRepository) *OperationLogService {
	return &OperationLogService{
		logRepo: logRepo,
	}
}

func (s *OperationLogService) LogOperation(
	c *gin.Context,
	userID uuid.UUID,
	operationType models.OperationType,
	resourceType models.ResourceType,
	resourceID *uuid.UUID,
	details interface{},
	result models.OperationResult,
	errorMessage string,
) error {
	var ipAddress string
	var userAgent string

	if c != nil {
		ipAddress = c.ClientIP()
		userAgent = c.Request.UserAgent()
	}

	var resourceIDStr *string
	if resourceID != nil {
		idStr := resourceID.String()
		resourceIDStr = &idStr
	}

	var detailsStr string
	if details != nil {
		if jsonBytes, err := json.Marshal(details); err == nil {
			detailsStr = string(jsonBytes)
		}
	}

	log := &models.OperationLog{
		UserID:       &userID,
		Operation:    operationType,
		ResourceType: resourceType,
		ResourceID:   resourceIDStr,
		Details:      detailsStr,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Result:       result,
		Error:        errorMessage,
	}

	if err := s.logRepo.Create(log); err != nil {
		return fmt.Errorf("failed to log operation: %w", err)
	}

	return nil
}

func (s *OperationLogService) GetLogs(filter models.OperationLogFilter) ([]models.OperationLog, int64, error) {
	logs, total, err := s.logRepo.FindAll(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, total, nil
}

func (s *OperationLogService) GetUserLogs(userID uuid.UUID, filter models.OperationLogFilter) ([]models.OperationLog, int64, error) {
	logs, total, err := s.logRepo.FindByUser(userID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user logs: %w", err)
	}
	return logs, total, nil
}

func (s *OperationLogService) GetSystemStats() (*models.SystemStats, error) {
	stats, err := s.logRepo.GetSystemStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get system stats: %w", err)
	}
	return stats, nil
}

func (s *OperationLogService) GetUserOperationStats(userID uuid.UUID, startDate, endDate time.Time) (map[string]int64, error) {
	stats, err := s.logRepo.GetUserOperationStats(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get user operation stats: %w", err)
	}
	return stats, nil
}

func (s *OperationLogService) CleanupOldLogs(days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	deletedCount, err := s.logRepo.DeleteOldLogs(cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old logs: %w", err)
	}
	return deletedCount, nil
}
