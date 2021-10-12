package gokamux

import (
	"fmt"

	"github.com/lovoo/goka"
)

type pipeline struct {
	modules []namedModule
	Verbose bool
}

type namedModule struct {
	name   string
	module Module
}

func newPipeline(modules ...string) (*pipeline, error) {
	p := pipeline{}

	for _, module := range modules {
		m, err := findModule(module)

		if err != nil {
			return nil, err
		}

		p.modules = append(p.modules, namedModule{name: module, module: m})
	}

	return &p, nil
}

type pipelineResult struct {
	Discard  bool
	Override bool
	Error    error
}

func (p pipeline) Run(ctx goka.Context, message interface{}, params []string) pipelineResult {
	r := pipelineResult{}

	for i, m := range p.modules {
		cbCtx := cbContext{GokaCtx: ctx}

		fmt.Printf("Executing module %d [%s]\n", i, m.name)

		m.module.Execute(&cbCtx, message, params)

		fmt.Printf("Module %d [%s]\n executed", i, m.name)

		if cbCtx.Error != nil {
			return pipelineResult{Error: cbCtx.Error}
		}

		if cbCtx.Discarded {
			return pipelineResult{Discard: true}
		}

		if cbCtx.OverridedMessage != nil {
			message = cbCtx.OverridedMessage
			r.Override = true
		}
	}

	return r
}
