package dto

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

type GetStatus struct {
	UserID        uuid.UUID
	Status        status.Status
	URL           string
	UploadedFiles int
}
