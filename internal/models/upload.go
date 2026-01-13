package models

import (
	"time"

	"github.com/google/uuid"
)

type UploadStatus string

const (
	UploadStatusPending   UploadStatus = "pending"
	UploadStatusUploading UploadStatus = "uploading"
	UploadStatusCompleted UploadStatus = "completed"
	UploadStatusFailed    UploadStatus = "failed"
	UploadStatusCanceled  UploadStatus = "canceled"
)

type UploadSession struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	FileName       string       `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize       int64        `gorm:"not null" json:"file_size"`
	FileHash       string       `gorm:"type:varchar(255)" json:"file_hash"`
	ParentID       *uuid.UUID   `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	ChunkSize      int64        `gorm:"not null" json:"chunk_size"`
	TotalChunks    int          `gorm:"not null" json:"total_chunks"`
	UploadedChunks int          `gorm:"default:0" json:"uploaded_chunks"`
	StoragePath    string       `gorm:"type:varchar(512)" json:"storage_path"`
	MimeType       string       `gorm:"type:varchar(100)" json:"mime_type"`
	Status         UploadStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ErrorMessage   string       `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiresAt      time.Time    `gorm:"index" json:"expires_at"`
}

func (UploadSession) TableName() string {
	return "upload_sessions"
}

type ChunkUploadRequest struct {
	UploadID   uuid.UUID `form:"upload_id" binding:"required"`
	ChunkIndex int       `form:"chunk_index" binding:"required,min=0"`
	ChunkSize  int64     `form:"chunk_size" binding:"required,min=1"`
	ChunkHash  string    `form:"chunk_hash" binding:"required"`
}

type ChunkUploadResponse struct {
	ChunkIndex      int     `json:"chunk_index"`
	Uploaded        bool    `json:"uploaded"`
	UploadedSize    int64   `json:"uploaded_size"`
	Progress        float64 `json:"progress"`
	CompletedChunks []int   `json:"completed_chunks"`
}

type UploadSessionResponse struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"user_id"`
	FileName        string       `json:"file_name"`
	FileSize        int64        `json:"file_size"`
	ParentID        *uuid.UUID   `json:"parent_id,omitempty"`
	ChunkSize       int64        `json:"chunk_size"`
	TotalChunks     int          `json:"total_chunks"`
	UploadedChunks  int          `json:"uploaded_chunks"`
	Progress        float64      `json:"progress"`
	Status          UploadStatus `json:"status"`
	CreatedAt       time.Time    `json:"created_at"`
	ExpiresAt       time.Time    `json:"expires_at"`
	CompletedChunks []int        `json:"completed_chunks"`
}

func (s *UploadSession) ToResponse(completedChunks []int) UploadSessionResponse {
	progress := float64(0)
	if s.TotalChunks > 0 {
		progress = float64(s.UploadedChunks) / float64(s.TotalChunks) * 100
	}

	return UploadSessionResponse{
		ID:              s.ID,
		UserID:          s.UserID,
		FileName:        s.FileName,
		FileSize:        s.FileSize,
		ParentID:        s.ParentID,
		ChunkSize:       s.ChunkSize,
		TotalChunks:     s.TotalChunks,
		UploadedChunks:  s.UploadedChunks,
		Progress:        progress,
		Status:          s.Status,
		CreatedAt:       s.CreatedAt,
		ExpiresAt:       s.ExpiresAt,
		CompletedChunks: completedChunks,
	}
}

type InitiateUploadRequest struct {
	FileName  string     `json:"file_name" binding:"required"`
	FileSize  int64      `json:"file_size" binding:"required,min=1"`
	FileHash  string     `json:"file_hash" binding:"required"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	ChunkSize int64      `json:"chunk_size" binding:"required,min=1"`
	MimeType  string     `json:"mime_type"`
}

type CompleteUploadRequest struct {
	UploadID uuid.UUID `json:"upload_id" binding:"required"`
}
