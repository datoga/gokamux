package gokamux

import (
	"fmt"
	"sync"
)

var (
	moduleMtx     sync.RWMutex
	moduleLoaders = make(map[string]ModuleLoader)
	modules       = make(map[string]Module)
)

func RegisterLoader(name string, moduleLoader ModuleLoader) {
	if moduleLoader == nil {
		panic(fmt.Errorf("nil module loader %s provided", name))
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	checkDupOrPanic(name)

	moduleLoaders[name] = moduleLoader
}

func RegisterModule(name string, module Module) {
	if module == nil {
		panic(fmt.Errorf("nil module %s provided", name))
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	checkDupOrPanic(name)

	modules[name] = module
}

func checkDupOrPanic(name string) {
	if _, dup := modules[name]; dup {
		panic(fmt.Errorf("module %s registered previously as module", name))
	}

	if _, dup := moduleLoaders[name]; dup {
		panic(fmt.Errorf("module loader %s registered previously as module loader", name))
	}
}

func instanceModule(name string, params ...string) (Module, error) {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	loader, found := moduleLoaders[name]

	if found {
		return loader.Init(params...), nil
	}

	module, found := modules[name]

	if !found {
		return nil, fmt.Errorf("module %s not found", name)
	}

	return module, nil
}

func mustInstanceModule(name string, params ...string) Module {
	m, err := instanceModule(name, params...)

	if err != nil {
		panic(err)
	}

	return m
}
