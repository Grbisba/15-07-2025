package inmemory

import (
	"bytes"
	"encoding/gob"
	"os"
)

type Option[K comparable, V any] interface {
	onStart(cache *Cache[K, V]) error
	onStop(cache *Cache[K, V]) error
}

type startOption[K comparable, V any] func(cache *Cache[K, V]) error

func (f startOption[K, V]) onStart(cache *Cache[K, V]) error {
	if cache == nil || f == nil {
		return ErrNilReference
	}
	cache.mu.Lock()
	err := f(cache)
	cache.mu.Unlock()
	return err
}

func (startOption[K, V]) onStop(*Cache[K, V]) error { return nil }

type stopOption[K comparable, V any] func(cache *Cache[K, V]) error

func (f stopOption[K, V]) onStop(cache *Cache[K, V]) error {
	if cache == nil || f == nil {
		return ErrNilReference
	}
	cache.mu.Lock()
	err := f(cache)
	cache.mu.Unlock()
	return err
}

func (stopOption[K, V]) onStart(*Cache[K, V]) error { return nil }

type fileSyncerOption[K comparable, V any] struct {
	filename string
	exist    bool
}

func fileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}

func WithSyncingWithFile[K comparable, V any](filename string) Option[K, V] {
	fso := fileSyncerOption[K, V]{
		filename: filename,
		exist:    fileExists(filename),
	}

	return fso
}

func (opt fileSyncerOption[K, V]) onStart(cache *Cache[K, V]) error {
	if cache == nil {
		return ErrNilReference
	}
	if !opt.exist {
		return nil
	}
	fi, err := os.Stat(opt.filename)
	if err != nil || fi == nil {
		return err
	}
	if fi.Size() == 0 {
		return nil
	}
	data, err := os.ReadFile(opt.filename)
	if err != nil {
		return err
	}

	cache.mu.Lock()
	table := cache.data
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&table)
	cache.data = table
	cache.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (opt fileSyncerOption[K, V]) onStop(cache *Cache[K, V]) error {
	if cache == nil {
		return ErrNilReference
	}
	if cache.data == nil || len(cache.data) == 0 {
		return nil
	}

	f, err := os.OpenFile(opt.filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	cache.mu.RLock()
	err = gob.NewEncoder(f).Encode(&cache.data)
	cache.mu.RUnlock()
	return err
}

func WithData[K comparable, V any](data map[K]V) Option[K, V] {
	if data == nil {
		return startOption[K, V](nil)
	}
	return startOption[K, V](func(cache *Cache[K, V]) error {
		cache.data = data
		return nil
	})
}
