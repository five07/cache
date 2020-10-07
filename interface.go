package cache

import "time"

// Cache interface
type Cache interface {
	Init()
	Keys() []string
	Has(key string) bool
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
	ExpiresIn(key string) time.Time // redis TTL <key>
}
