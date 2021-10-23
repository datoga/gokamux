package modules

import (
	"fmt"
	"sync"

	"github.com/datoga/gokamux/modules/model"
)

var (
	moduleMtx sync.RWMutex
	registry  = make(map[string]model.Module)
)

func Register(name string, module model.Module) error {
	if name == "" {
		return fmt.Errorf("an unique name must be provided for the module")
	}

	if module == nil {
		return fmt.Errorf("nil module %s provided", name)
	}

	moduleMtx.Lock()
	defer moduleMtx.Unlock()

	if _, dup := registry[name]; dup {
		return fmt.Errorf("module %s registered previously", name)
	}

	registry[name] = module

	return nil
}

func MustRegister(name string, module model.Module) {
	if err := Register(name, module); err != nil {
		panic(err)
	}
}

func List() []string {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	var mNames []string

	for k := range registry {
		mNames = append(mNames, k)
	}

	return mNames
}

func Instance(id string, params ...string) (model.Module, error) {
	moduleMtx.RLock()
	defer moduleMtx.RUnlock()

	module, found := registry[id]

	if !found {
		return nil, fmt.Errorf("module %s not found", id)
	}

	if err := module.Init(params...); err != nil {
		return nil, fmt.Errorf("failed on init for module %s with error %v", id, err)
	}

	return module, nil
}
