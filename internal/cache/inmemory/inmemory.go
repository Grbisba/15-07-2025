package inmemory

import (
	"sync"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Grbisba/15-07-2025/internal/cache"
)

var _ cache.Cache[int, int] = (*Cache[int, int])(nil)

type Cache[K comparable, V any] struct {
	data    map[K]V
	mu      *sync.RWMutex
	log     *zap.Logger
	options []Option[K, V]
}

func New[K comparable, V any](logger *zap.Logger, opts ...Option[K, V]) (*Cache[K, V], error) {
	if logger == nil {
		logger = zap.NewNop().Named("cache")
	}

	c := &Cache[K, V]{
		data:    make(map[K]V, 100),
		mu:      &sync.RWMutex{},
		log:     logger.Named("cache"),
		options: opts,
	}

	c.log.Info("initializing cache")

	return c, c.onStart()
}

func (c *Cache[K, V]) Close() error {
	return c.onStop()
}

func (c *Cache[K, V]) onStart() error {
	var err error

	c.log.Info("applying options on cache")

	for _, opt := range c.options {
		err = multierr.Append(err, opt.onStart(c))
	}

	return err
}

func (c *Cache[K, V]) onStop() error {
	var err error

	c.log.Info("stopping cache")

	for _, opt := range c.options {
		err = multierr.Append(err, opt.onStop(c))
	}

	return err
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.log.Debug("setting object to cache", zap.Any("key", key))
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.log.Debug("getting object from cache", zap.Any("key", key))

	c.mu.RLock()
	value, ok = c.data[key]
	c.mu.RUnlock()

	return value, ok
}
