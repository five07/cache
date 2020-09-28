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
	store      = make(map[string]Cache)
	registry   = make(map[string]Cache)
	registryMu sync.RWMutex
)

// Register func
func Register(name string, adapter Cache) {
	registryMu.Lock()
	defer registryMu.Unlock()
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
		registryMu.RLock()
		adapter, ok := registry[lib.Config.AdapterName]
		registryMu.RUnlock()
		if !ok {
			panic("cache: unknown adapter " + lib.Config.AdapterName + " (forgotten import?)")
		}

		cache = adapter.(Cache)
		cache.Init()
		store[key] = cache
	}

	return cache
}
