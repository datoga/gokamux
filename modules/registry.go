package modules

import (
	"fmt"
	"sync"
)

type ModuleDefinition struct {
	Name       string
	Module     Module
	Configurer Configurer
}

var (
	moduleMtx sync.RWMutex
	modules   = make(map[string]ModuleDefinition)
)

func Register(id string, module ModuleDefinition) error {
	if id == "" {
		return fmt.Errorf("an unique id must be provided for the module")
	}

	if module.Module == nil {
		return fmt.Errorf("nil module %s provided", id)
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	if _, dup := modules[id]; dup {
		return fmt.Errorf("module %s registered previously", id)
	}

	modules[id] = module

	return nil
}

func List() []ModuleDefinition {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	var mDefs []ModuleDefinition

	for _, m := range modules {
		mDefs = append(mDefs, m)
	}

	return mDefs
}

func Instance(id string, params ...string) (Module, error) {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	module, found := modules[id]

	if !found {
		return nil, fmt.Errorf("module %s not found", id)
	}

	if len(params) > 0 {
		if module.Configurer == nil {
			fmt.Println("no found configurer with params, it will register the module but params will be ignored")
		} else {
			err := module.Configurer.Configure(params...)

			if err != nil {
				return nil, err
			}
		}
	}

	return module.Module, nil
}
