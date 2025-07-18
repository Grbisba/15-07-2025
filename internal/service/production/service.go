package production

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/cache"
	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/encryptor"
	"github.com/Grbisba/15-07-2025/internal/pkg/manager"
	"github.com/Grbisba/15-07-2025/internal/pkg/storage"
	"github.com/Grbisba/15-07-2025/internal/pkg/token"
	"github.com/Grbisba/15-07-2025/internal/service"
)

var _ service.Service = (*Service)(nil)

type Service struct {
	log   *zap.Logger
	cfg   *config.Config
	tasks cache.Cache[uuid.UUID, dto.Task]
	store storage.UserStorage
	dm    manager.Manager
	enc   encryptor.Interface
	prv   token.Provider
}

func New(
	log *zap.Logger,
	cfg *config.Config,
	tasks cache.Cache[uuid.UUID, dto.Task],
	store storage.UserStorage,
	dm manager.Manager,
	ecn encryptor.Interface,
	prv token.Provider,
) (*Service, error) {
	if log == nil {
		log = zap.NewNop()
		log.Named("service")

	}
	s := &Service{
		log:   log.Named("service"),
		cfg:   cfg,
		tasks: tasks,
		store: store,
		dm:    dm,
		enc:   ecn,
		prv:   prv,
	}

	if err := s.validate(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) validate() error {
	switch {
	case s.cfg == nil:
		log.Error(
			"cfg unexpectedly nil",
			zap.Any("cfg", s.cfg),
		)
		return ErrNilConfig
	case s.store == nil:
		log.Error(
			"user-store unexpectedly nil",
			zap.Any("store", s.store),
		)
		return ErrNilStore
	case s.tasks == nil:
		log.Error(
			"cache unexpectedly nil",
			zap.Any("tasks", s.tasks),
		)
		return ErrNilCache
	case s.prv == nil:
		log.Error(
			"provider unexpectedly nil",
			zap.Any("prv", s.prv),
		)
		return ErrNilProvider
	case s.enc == nil:
		log.Error(
			"encryptor unexpectedly nil",
			zap.Any("enc", s.enc),
		)
		return ErrNilEncryptor
	case s.dm == nil:
		log.Error(
			"manager unexpectedly nil",
			zap.Any("dm", s.dm),
		)
		return ErrNilManager
	default:
		s.log.Info(
			"service successfully validated",
			zap.Any("service", s),
		)
		return nil
	}
}

func (s *Service) Shutdown() error {
	return multierr.Combine(s.tasks.Close(), s.store.Close())
}
