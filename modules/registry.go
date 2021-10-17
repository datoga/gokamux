package modules

import (
	"fmt"
	"sync"

	"github.com/datoga/gokamux/modules/model"
)

type ModuleDefinition struct {
	Name   string
	Module model.Module
}

var (
	moduleMtx sync.RWMutex
	registry  = make(map[string]ModuleDefinition)
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

	if _, dup := registry[id]; dup {
		return fmt.Errorf("module %s registered previously", id)
	}

	registry[id] = module

	return nil
}

func List() []ModuleDefinition {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	var mDefs []ModuleDefinition

	for _, m := range registry {
		mDefs = append(mDefs, m)
	}

	return mDefs
}

func Instance(id string, params ...string) (model.Module, error) {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	module, found := registry[id]

	if !found {
		return nil, fmt.Errorf("module %s not found", id)
	}

	if err := module.Module.Init(params...); err != nil {
		return nil, fmt.Errorf("failed on init for module %s with error %v", id, err)
	}

	return module.Module, nil
}
