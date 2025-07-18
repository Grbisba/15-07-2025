package dto

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

type (
	Task struct {
		UserID uuid.UUID
		ID     uuid.UUID
		Status status.Status
		Files  []File
	}
)
