package token

import (
	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

type Provider interface {
	CreateAccessTokenForUser(userID uuid.UUID) (string, error)
	CreateRefreshTokenForUser(userID uuid.UUID) (string, error)
	GetDataFromToken(raw string) (*dto.UserDataInToken, error)
	CreateAccessAndRefreshTokenForUser(userID uuid.UUID) (string, string, error)
}
