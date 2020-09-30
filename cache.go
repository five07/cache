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
	mutex    sync.RWMutex
	store    = make(map[string]Cache)
	registry = make(map[string]Cache)
)

// Register func
func Register(name string, adapter Cache) {
	mutex.Lock()
	defer mutex.Unlock()
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
		mutex.RLock()
		adapter, ok := registry[lib.Config.AdapterName]
		mutex.RUnlock()
		if !ok {
			panic("cache: unknown adapter " + lib.Config.AdapterName + " (forgotten import?)")
		}

		mutex.Lock()
		defer mutex.Unlock()

		cache = adapter.(Cache)
		cache.Init()
		store[key] = cache
	}

	return cache
}
