package main

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
)

var message string

func Configure(params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("no params provided")
	}

	message = params[0]

	return nil
}

func Process(ctx modules.Context, msg string) {
	fmt.Println("Processor changer with message", message)

	ctx.OverrideMessage(message)
}
