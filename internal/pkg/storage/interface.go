package storage

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

type UserStorage interface {
	Get(k string) (dto.User, bool)
	Set(login string, id uuid.UUID, v dto.User)
	Close() error
}
