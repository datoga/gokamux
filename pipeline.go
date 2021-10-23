package gokamux

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/model"
	"github.com/lovoo/goka"
)

type compiledStep struct {
	ID             string
	ModuleName     string
	ModuleInstance model.Module
}

type pipeline struct {
	steps         []Step
	compiledSteps []compiledStep
	Verbose       bool
}

func newPipeline(steps ...Step) *pipeline {
	return &pipeline{
		steps: steps,
	}
}

type pipelineRunner struct {
	compiledSteps []compiledStep
}

func (p *pipeline) Compile() (*pipelineRunner, error) {
	var pSteps []compiledStep

	for _, step := range p.steps {
		module, err := modules.Instance(step.Module, step.Params...)

		if err != nil {
			return nil, fmt.Errorf("failed looking for module %s in registry with error %v", step.Module, err)
		}

		pStep := compiledStep{
			ID:             step.ID,
			ModuleName:     step.Module,
			ModuleInstance: module,
		}

		pSteps = append(pSteps, pStep)
	}

	return &pipelineRunner{
		compiledSteps: pSteps,
	}, nil
}

type pipelineResult struct {
	Discard  bool
	Override bool
	Error    error
}

func (p pipelineRunner) Run(ctx goka.Context, message *string) pipelineResult {
	r := pipelineResult{}

	for i, m := range p.compiledSteps {
		cbCtx := cbContext{GokaCtx: ctx}

		fmt.Printf("Executing step %d [%s] with module %s\n", i, m.ID, m.ModuleName)

		if err := m.ModuleInstance.Process(&cbCtx, *message); err != nil {
			return pipelineResult{Error: err}
		}

		fmt.Printf("Step %d [%s]\n executed successfully", i, m.ID)

		if cbCtx.Error != nil {
			return pipelineResult{Error: cbCtx.Error}
		}

		if cbCtx.Discarded {
			return pipelineResult{Discard: true}
		}

		if cbCtx.OverridedMessage != "" {
			message = &cbCtx.OverridedMessage
			r.Override = true
		}
	}

	return r
}
