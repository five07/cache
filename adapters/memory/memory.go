package memory

import (
	"reflect"
	"sync"
	"time"

	"github.com/five07/cache"
)

// AdapterName must be unique
const AdapterName = "memory"
const defaultTTL = 3600

var mutex sync.RWMutex

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
	mutex.RLock()
	keys := reflect.ValueOf(a.store).MapKeys()
	mutex.RUnlock()

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k.String())
	}

	return result
}

func (a *Memory) exists(key string) bool {
	mutex.RLock()
	_, exists := a.store[key]
	mutex.RUnlock()
	return exists
}

// Has func
func (a *Memory) Has(key string) bool {
	return a.ExpiresIn(key) > 0
}

// Get func
func (a *Memory) Get(key string) interface{} {
	if a.Has(key) && a.ExpiresIn(key) > 0 {
		mutex.RLock()
		defer mutex.RUnlock()
		return a.store[key].data
	}

	return nil
}

// Set func
func (a *Memory) Set(key string, value interface{}, ttl uint64) {
	if ttl == 0 {
		ttl = defaultTTL
	}

	mutex.Lock()
	a.store[key] = CacheValue{value, time.Now().Add(time.Second * time.Duration(ttl))}
	mutex.Unlock()
}

// Delete func
func (a *Memory) Delete(key string) {
	mutex.Lock()
	delete(a.store, key)
	mutex.Unlock()
}

// Clear func
func (a *Memory) Clear() {
	mutex.Lock()
	a.store = make(map[string]CacheValue)
	mutex.Unlock()
}

// ExpiresIn func
func (a *Memory) ExpiresIn(key string) uint64 {
	if a.exists(key) {
		mutex.RLock()
		d := a.store[key].expiration.Sub(time.Now()).Seconds()
		mutex.RUnlock()

		if d > 0 {
			return uint64(d)
		}

		mutex.Lock()
		a.Delete(key)
		mutex.Unlock()
	}

	return 0
}
