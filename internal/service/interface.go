package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

type Service interface {
	Shutdown() error

	CreateTask(ctx context.Context, userID uuid.UUID) (*dto.Task, error)
	GetStatus(_ context.Context, id uuid.UUID, userID uuid.UUID) (*dto.GetStatus, error)
	UploadFile(ctx context.Context, req dto.UploadFile) (*dto.Task, error)
	Download(ctx context.Context, id uuid.UUID) (*dto.Download, error)

	CreateUser(ctx context.Context, req dto.CreateUser) (*dto.Token, error)
	Login(ctx context.Context, req dto.Login) (*dto.Token, error)
	RefreshToken(ctx context.Context, token string) (*dto.Token, error)
	CheckUserAuth(_ context.Context, token string) (*dto.UserDataInToken, error)
}
