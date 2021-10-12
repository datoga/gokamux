package gokamux

import (
	"fmt"
	"sync"
)

var (
	modulesMtx sync.RWMutex
	modules    = make(map[string]Module)
)

func Register(name string, module Module) {
	if module == nil {
		panic(fmt.Errorf("nil module %s provided", name))
	}

	modulesMtx.Lock()
	defer modulesMtx.Unlock()

	if _, dup := modules[name]; dup {
		panic(fmt.Errorf("module %s registered previously", name))
	}

	modules[name] = module
}

func findModule(name string) (Module, error) {
	modulesMtx.RLock()
	defer modulesMtx.RUnlock()
	
	m, found := modules[name]

	if !found {
		return nil, fmt.Errorf("module %s not found", name)
	}

	return m, nil
}
