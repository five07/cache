package memory

import (
	"reflect"
	"time"

	"github.com/five07/cache"
)

// AdapterName must be unique
const AdapterName = "memory"
const defaultTTL = 3600

// CacheValue type
type CacheValue struct {
	data       interface{}
	expiration time.Time
}

// Memory cache type
type Memory struct {
	Config interface{}
	store  map[string]CacheValue
}

func init() {
	cache.Register("memory", &Memory{})
}

// Init func
func (a *Memory) Init() {
	a.store = make(map[string]CacheValue)
}

// Keys func
func (a *Memory) Keys() []string {
	keys := reflect.ValueOf(a.store).MapKeys()
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k.String())
	}

	return result
}

func (a *Memory) exists(key string) bool {
	_, exists := a.store[key]
	return exists
}

// Has func
func (a *Memory) Has(key string) bool {
	return a.ExpiresIn(key) > 0
}

// Get func
func (a *Memory) Get(key string) interface{} {
	if a.Has(key) && a.ExpiresIn(key) > 0 {
		return a.store[key].data
	}

	return nil
}

// Set func
func (a *Memory) Set(key string, value interface{}, ttl uint64) {
	if ttl == 0 {
		ttl = defaultTTL
	}

	a.store[key] = CacheValue{value, time.Now().Add(time.Second * time.Duration(ttl))}
}

// Delete func
func (a *Memory) Delete(key string) {
	delete(a.store, key)
}

// Clear func
func (a *Memory) Clear() {
	a.store = make(map[string]CacheValue)
}

// ExpiresIn func
func (a *Memory) ExpiresIn(key string) uint64 {
	if a.exists(key) {
		d := a.store[key].expiration.Sub(time.Now()).Seconds()

		if d > 0 {
			return uint64(d)
		}

		a.Delete(key)
	}

	return 0
}
