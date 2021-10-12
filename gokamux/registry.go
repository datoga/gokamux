package gokamux

import (
	"fmt"
	"sync"
)

var (
	moduleLoaderMtx sync.RWMutex
	moduleLoaders   = make(map[string]ModuleLoader)
)

func Register(name string, moduleLoader ModuleLoader) {
	if moduleLoader == nil {
		panic(fmt.Errorf("nil module loader %s provided", name))
	}

	moduleLoaderMtx.Lock()
	defer moduleLoaderMtx.Unlock()

	if _, dup := moduleLoaders[name]; dup {
		panic(fmt.Errorf("module loader %s registered previously", name))
	}

	moduleLoaders[name] = moduleLoader
}

func findModuleLoader(name string) (ModuleLoader, error) {
	moduleLoaderMtx.RLock()
	defer moduleLoaderMtx.RUnlock()

	m, found := moduleLoaders[name]

	if !found {
		return nil, fmt.Errorf("module %s not found", name)
	}

	return m, nil
}

func mustModule(name string) Module {
	m, err := findModule(name)

	if err != nil {
		panic(err)
	}

	return m
}
