package main

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
)

func Init(params ...string) (modules.Module, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("no params provided")
	}

	return New(params[0]), nil
}

type Processor struct {
	message string
}

func New(message string) *Processor {
	return &Processor{message: message}
}

func (p Processor) Process(ctx modules.Context, msg string) {
	fmt.Println("Processor changer with message", p.message)

	ctx.OverrideMessage(p.message)
}
