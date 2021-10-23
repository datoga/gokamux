package main

import (
	"fmt"

	"github.com/datoga/gokamux/modules/model"
)

func Load() model.Module {
	return &sampleParams{}
}

type sampleParams struct {
	message string
}

func (s *sampleParams) Init(params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("no params provided")
	}

	s.message = params[0]

	return nil
}

func (s sampleParams) Process(ctx model.Context, msg string) error {
	fmt.Println("Processor changer with message", s.message)

	ctx.OverrideMessage(s.message)

	return nil
}
