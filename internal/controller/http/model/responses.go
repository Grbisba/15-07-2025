package model

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

type (
	TaskResponse struct {
		ID     uuid.UUID     `json:"id"`
		Status status.Status `json:"status"`
	}
	GetStatusResponse struct {
		UserID        uuid.UUID     `json:"user_id"`
		Status        status.Status `json:"status"`
		URL           string        `json:"url"`
		UploadedFiles int           `json:"uploaded_files"`
	}
	CreateUserResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	RefreshTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
