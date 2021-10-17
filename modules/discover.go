package modules

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/datoga/gokamux/modules/model"
)

func Discover(path string) ([]ModuleDefinition, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("path %s not readable with error %v", path, err)
	}

	var modules []ModuleDefinition

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		module, err := discoverModule(filepath.Join(path, f.Name()))

		if err != nil {
			return nil, err
		}

		modules = append(modules, *module)
	}

	return modules, nil
}

func discoverModule(name string) (*ModuleDefinition, error) {
	plug, err := plugin.Open(name)

	if err != nil {
		return nil, fmt.Errorf("failed opening plugin %s with error %v", name, err)
	}

	processPlugin, err := plug.Lookup("Load")

	if err != nil {
		return nil, fmt.Errorf("failed looking symbols on plugin %s with error %v", name, err)
	}

	log.Printf("%T, %+v", processPlugin, processPlugin)

	fnLoad, ok := processPlugin.(func() model.Module)

	if !ok {
		return nil, fmt.Errorf("failed taking function process for plugin %s", name)
	}

	pluginName := name

	if idx := strings.Index(name, "."); idx != -1 {
		pluginName = name[:idx]
	}

	return &ModuleDefinition{
		Name:   pluginName,
		Module: fnLoad(),
	}, nil
}

func RegisterModules(modules ...ModuleDefinition) error {
	for _, module := range modules {
		if err := Register(module.Name, module); err != nil {
			return err
		}
	}

	return nil
}

func DiscoverAndRegister(path string) error {
	modules, err := Discover(path)

	if err != nil {
		return fmt.Errorf("failed discovering modules with error %v", err)
	}

	err = RegisterModules(modules...)

	if err != nil {
		return fmt.Errorf("failed registering modules with error %v", err)
	}

	return nil
}

func MustDiscoverAndRegister(path string) {
	if err := DiscoverAndRegister(path); err != nil {
		panic(err)
	}
}
