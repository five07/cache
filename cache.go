package cache

import (
	"sync"
)

// Library type
type Library struct {
	Config *Config
}

// Config type
type Config struct {
	AdapterName string
}

var (
	lock     sync.RWMutex
	store    = make(map[string]Cache)
	registry = make(map[string]Cache)
)

// Register func
func Register(name string, adapter Cache) {
	lock.Lock()
	defer lock.Unlock()
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, dup := registry[name]; dup {
		panic("cache: Register called twice for adapter " + name)
	}
	registry[name] = adapter
}

// Get func
func (lib Library) Get(key string) Cache {
	var cache Cache
	var ok bool

	cache, ok = store[key]

	if !ok {
		lock.RLock()
		adapter, ok := registry[lib.Config.AdapterName]
		lock.RUnlock()
		if !ok {
			panic("cache: unknown adapter " + lib.Config.AdapterName + " (forgotten import?)")
		}

		cache = adapter.(Cache)
		cache.Init()
		store[key] = cache
	}

	return cache
}
