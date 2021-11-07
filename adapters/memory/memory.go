package memory

import (
	"reflect"
	"sync"
	"time"

	"github.com/five07/cache"
)

// AdapterName must be unique
const AdapterName = "memory"

// Item type
type Item struct {
	data       interface{}
	expiration int64
}

// Memory cache type
type Memory struct {
	Config interface{}
	store  map[string]Item
	mutex  sync.RWMutex
}

func init() {
	cache.Register("memory", &Memory{})
}

// Init func
func (a *Memory) Init() {
	a.store = make(map[string]Item)
}

// Keys func
func (a *Memory) Keys() []string {
	a.mutex.RLock()
	keys := reflect.ValueOf(a.store).MapKeys()
	a.mutex.RUnlock()

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k.String())
	}

	return result
}

// Has func
func (a *Memory) Has(key string) bool {
	_, exists := a.Get(key)
	return exists
}

// Get func
func (a *Memory) Get(key string) (interface{}, bool) {
	a.mutex.RLock()

	item, exists := a.store[key]
	if !exists {
		a.mutex.RUnlock()
		return nil, false
	}

	if a.expired(&item) {
		a.mutex.RUnlock()
		return nil, false
	}

	a.mutex.RUnlock()
	return item.data, true
}

// Set func
func (a *Memory) Set(key string, value interface{}, ttl time.Duration) {
	var exp int64

	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	a.mutex.Lock()
	a.store[key] = Item{value, exp}
	a.mutex.Unlock()
}

// Delete func
func (a *Memory) Delete(key string) {
	a.mutex.Lock()
	delete(a.store, key)
	a.mutex.Unlock()
}

// Clear func
func (a *Memory) Clear() {
	a.mutex.Lock()
	a.store = make(map[string]Item)
	a.mutex.Unlock()
}

// ExpiresIn func
func (a *Memory) ExpiresIn(key string) time.Time {
	item, exists := a.Get(key)
	if !exists {
		return time.Time{}
	}

	return time.Unix(0, item.expiration)
}

func (a *Memory) expired(item *Item) bool {
	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			return true
		}
	}

	return false
}
