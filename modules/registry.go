package modules

import (
	"fmt"
	"sync"
)

var (
	moduleMtx         sync.RWMutex
	modulesLoader     = make(map[string]ModuleLoader)
	modulesStandalone = make(map[string]Module)
)

func RegisterLoader(name string, moduleLoader ModuleLoader) {
	if moduleLoader == nil {
		panic(fmt.Errorf("nil module loader %s provided", name))
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	checkDupOrPanic(name)

	modulesLoader[name] = moduleLoader
}

func RegisterModule(name string, module Module) {
	if module == nil {
		panic(fmt.Errorf("nil module %s provided", name))
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	checkDupOrPanic(name)

	modulesStandalone[name] = module
}

func checkDupOrPanic(name string) {
	if _, dup := modulesStandalone[name]; dup {
		panic(fmt.Errorf("module %s registered previously as module", name))
	}

	if _, dup := modulesLoader[name]; dup {
		panic(fmt.Errorf("module loader %s registered previously as module loader", name))
	}
}

func InstanceModule(name string, params ...string) (Module, error) {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	loader, found := modulesLoader[name]

	if found {
		return loader.Init(params...)
	}

	module, found := modulesStandalone[name]

	if !found {
		return nil, fmt.Errorf("module %s not found", name)
	}

	if len(params) > 0 {
		fmt.Println("no found loader with params, it will instance a standalone module but params will be ignored")
	}

	return module, nil
}
