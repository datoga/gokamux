package modules

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/datoga/gokamux/modules/model"
)

type moduleDefinition struct {
	Name   string
	Module model.Module
}

func discover(path string) ([]moduleDefinition, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("path %s not readable with error %v", path, err)
	}

	var modules []moduleDefinition

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		module, err := discoverModule(filepath.Join(path, f.Name()))

		if err != nil {
			return nil, err
		}

		pluginName := f.Name()

		if idx := strings.Index(f.Name(), "."); idx != -1 {
			pluginName = f.Name()[:idx]
		}

		modules = append(modules, moduleDefinition{Name: pluginName, Module: *module})
	}

	return modules, nil
}

func discoverModule(name string) (*model.Module, error) {
	plug, err := plugin.Open(name)

	if err != nil {
		return nil, fmt.Errorf("failed opening plugin %s with error %v", name, err)
	}

	processPlugin, err := plug.Lookup("Load")

	if err != nil {
		return nil, fmt.Errorf("failed looking symbols on plugin %s with error %v", name, err)
	}

	fnLoad, ok := processPlugin.(func() model.Module)

	if !ok {
		return nil, fmt.Errorf("failed taking function process for plugin %s", name)
	}

	module := fnLoad()

	return &module, nil
}

func registerModules(modules ...moduleDefinition) error {
	for _, module := range modules {
		if err := Register(module.Name, module.Module); err != nil {
			return err
		}
	}

	return nil
}

func DiscoverAndRegister(path string) error {
	modules, err := discover(path)

	if err != nil {
		return fmt.Errorf("failed discovering modules with error %v", err)
	}

	err = registerModules(modules...)

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
