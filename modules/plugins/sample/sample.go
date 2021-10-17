package main

import (
	"fmt"

	"github.com/datoga/gokamux/modules/model"
)

func Load() model.Module {
	return sample{}
}

type sample struct{}

func (sample) Init(params ...string) error {
	return nil
}

func (sample) Process(ctx model.Context, msg string) {
	gokaCtx := ctx.GokaContext()

	fmt.Printf(
		"Message received, timestamp = %v, topic = %s, offset = %d, partition = %d, key = %s, bytes = %d, headers = %v, message = %s\n",
		gokaCtx.Timestamp(),
		gokaCtx.Topic(),
		gokaCtx.Offset(),
		gokaCtx.Partition(),
		gokaCtx.Key(),
		len(msg),
		gokaCtx.Headers(),
		msg,
	)
}
