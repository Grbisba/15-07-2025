package main

import (
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/Grbisba/loggr"

	"github.com/Grbisba/15-07-2025/internal/cache"
	"github.com/Grbisba/15-07-2025/internal/cache/inmemory"
	"github.com/Grbisba/15-07-2025/internal/config"
	"github.com/Grbisba/15-07-2025/internal/controller"
	"github.com/Grbisba/15-07-2025/internal/controller/http/fiber"
	"github.com/Grbisba/15-07-2025/internal/model/dto"
	"github.com/Grbisba/15-07-2025/internal/pkg/encryptor"
	"github.com/Grbisba/15-07-2025/internal/pkg/encryptor/hash"
	"github.com/Grbisba/15-07-2025/internal/pkg/manager"
	"github.com/Grbisba/15-07-2025/internal/pkg/manager/builder"
	"github.com/Grbisba/15-07-2025/internal/pkg/storage"
	"github.com/Grbisba/15-07-2025/internal/pkg/storage/user"
	"github.com/Grbisba/15-07-2025/internal/pkg/token"
	"github.com/Grbisba/15-07-2025/internal/pkg/token/jwt"
	"github.com/Grbisba/15-07-2025/internal/service"
	"github.com/Grbisba/15-07-2025/internal/service/production"
)

func main() {
	fx.New(buildOptions()).Run()
}

func buildOptions() fx.Option {
	return fx.Options(
		fx.WithLogger(zapLogger),
		fx.Provide(
			newLogger,
			config.New,

			fx.Annotate(builder.NewDirManager, fx.As(new(manager.Manager))),
			fx.Annotate(production.New, fx.As(new(service.Service))),
			fx.Annotate(newTasksCache, fx.As(new(cache.Cache[uuid.UUID, dto.Task]))),
			fx.Annotate(newUsersStorage, fx.As(new(storage.UserStorage))),
			fx.Annotate(hash.New, fx.As(new(encryptor.Interface))),
			fx.Annotate(jwt.NewProvider, fx.As(new(token.Provider))),
			controller.AsController(fiber.New),
		),
		fx.Invoke(
			controller.RunControllersFx,
		),
	)
}

func newLogger() *zap.Logger {
	l, _ := loggr.New()
	return l.Named("downloader")
}

func zapLogger(log *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{
		Logger: log.Named("fx"),
	}
}

func newTasksCache(log *zap.Logger, cfg *config.Config) (*inmemory.Cache[uuid.UUID, dto.Task], error) {
	im, err := inmemory.New[uuid.UUID, dto.Task](
		log,
		inmemory.WithSyncingWithFile[uuid.UUID, dto.Task](cfg.Cache.TaskCacheDataFile),
	)
	if err != nil {
		log.Error(
			"got error while initializing in-memory cache",
			zap.Error(err),
		)

		return nil, err
	}
	log.Info("cache successfully initialized")

	return im, nil
}

func newUsersStorage(log *zap.Logger, cfg *config.Config) (*user.Storage, error) {
	lim, err := inmemory.New[string, dto.User](
		log,
		inmemory.WithSyncingWithFile[string, dto.User](cfg.Cache.LoginUserCacheDataFile),
	)
	if err != nil {
		log.Error(
			"got error while initializing in-memory cache",
			zap.Error(err),
		)

		return nil, err
	}

	iim, err := inmemory.New[uuid.UUID, dto.User](
		log,
		inmemory.WithSyncingWithFile[uuid.UUID, dto.User](cfg.Cache.IDUserCacheDataFile),
	)
	if err != nil {
		log.Error(
			"got error while initializing in-memory cache",
			zap.Error(err),
		)

		return nil, err
	}

	log.Info("cache successfully initialized")

	return user.New(iim, lim), nil
}
