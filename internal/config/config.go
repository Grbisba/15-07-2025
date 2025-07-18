package config

import (
	"context"
	"os"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var configPath = os.Getenv("CONFIG_FILE_PATH")

var defaultConfig = &Config{}

type Config struct {
	Controller *Controller `config:"controller" toml:"controller" yaml:"controller" json:"controller"`
	Cache      *Cache      `config:"cache" toml:"cache" yaml:"cache" json:"cache"`
	Loader     *Loader     `config:"loader" toml:"loader" yaml:"loader" json:"loader"`
	JWT        *Token      `config:"jwt" toml:"jwt" yaml:"jwt" json:"jwt"`
}

func (c *Config) copy() *Config {
	if c == nil {
		return nil
	}

	return &Config{
		Controller: c.Controller.copy(),
		Cache:      c.Cache.copy(),
		Loader:     c.Loader.copy(),
	}
}

func New(log *zap.Logger) (*Config, error) {
	cfg := defaultConfig.copy()

	if configPath == "" {
		log.Error(
			"config file path not specified",
			zap.String("config-path", configPath),
		)
		return nil, errors.New("config file path is required")
	}

	l := confita.NewLoader(
		file.NewBackend(configPath),
	)

	err := l.Load(context.Background(), cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error while loading config")
	}

	log.Named("config").Info("loaded config", zap.Any("config", cfg))

	return cfg, nil
}
