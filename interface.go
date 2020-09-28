package cache

// Cache interface
type Cache interface {
	Init()
	Keys() []string
	Has(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{}, ttl uint64)
	Delete(key string)
	Clear()
	ExpiresIn(key string) uint64 // redis TTL <key>
}
