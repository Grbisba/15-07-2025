package cache

//go:generate mockgen -source=interface.go -destination=../../mocks/mockCache/mock.go -package=mockCache

type Cache[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (value V, ok bool)
	Close() error
}
