package modules

type fnProcess func(ctx Context, msg string)

type helperProcess struct {
	fnProcess
}

func (h helperProcess) Process(ctx Context, msg string) {
	h.fnProcess(ctx, msg)
}

type fnConfigure func(params ...string) error

type helperConfigure struct {
	fnConfigure
}

func (h helperConfigure) Configure(params ...string) error {
	return h.fnConfigure(params...)
}
