package log

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/model"
)

func init() {
	modules.MustRegister("log", &log{})
}

type log struct{}

func (log) Init(params ...string) error {
	return nil
}

func (log) Process(ctx model.Context, msg string) error {
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

	return nil
}
