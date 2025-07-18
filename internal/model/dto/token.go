package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDataInToken struct {
	ID        uuid.UUID
	IsAccess  bool
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
}

type AuthData struct {
	Login *string
	ID    *uuid.UUID
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
