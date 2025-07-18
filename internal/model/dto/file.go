package dto

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/pkg/status"
)

type (
	File struct {
		Name   string
		URL    string
		Ext    string
		Status status.Status
	}
	UploadFile struct {
		UserID uuid.UUID
		TaskID uuid.UUID
		URLs   []string
	}
	Download struct {
		UserID   uuid.UUID
		RawZip   []byte
		Status   status.Status
		Filename string
	}
)
