package caching

type Cache interface {
	Set(key, value interface{})

	Get(key interface{}) (interface{}, bool)

	Contains(key interface{}) bool
}
