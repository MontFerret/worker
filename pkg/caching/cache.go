package caching

type Cache[T any] interface {
	Set(key string, value T)

	Get(key string) (T, bool)

	Contains(key string) bool
}
