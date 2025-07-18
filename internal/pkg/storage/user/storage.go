package user

import (
	"github.com/google/uuid"
	"go.uber.org/multierr"

	"github.com/Grbisba/15-07-2025/internal/cache/inmemory"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
)

type Storage struct {
	IDStorage    *inmemory.Cache[uuid.UUID, dto.User]
	LoginStorage *inmemory.Cache[string, dto.User]
}

func New(
	IDS *inmemory.Cache[uuid.UUID, dto.User],
	LS *inmemory.Cache[string, dto.User],
) *Storage {
	return &Storage{
		IDStorage:    IDS,
		LoginStorage: LS,
	}
}

func (s *Storage) Get(k string) (dto.User, bool) {
	id, err := uuid.Parse(k)
	if err != nil {
		return s.getByLogin(k)
	}

	return s.getByID(id)
}

func (s *Storage) getByID(k uuid.UUID) (dto.User, bool) {
	return s.IDStorage.Get(k)
}

func (s *Storage) getByLogin(k string) (dto.User, bool) {
	return s.LoginStorage.Get(k)
}

func (s *Storage) Set(login string, id uuid.UUID, v dto.User) {
	s.LoginStorage.Set(login, v)
	s.IDStorage.Set(id, v)
	return
}

func (s *Storage) Close() error {
	return multierr.Combine(s.IDStorage.Close(), s.LoginStorage.Close())
}
