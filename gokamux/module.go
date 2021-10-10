package gokamux

import "github.com/lovoo/goka"

type ModuleFn func(ctx Context, msg interface{}, params []string)

type Module struct {
	Name string
	fn   ModuleFn
}

type FilterResult bool

var Allow = FilterResult(false)
var Discard = FilterResult(true)

func New(name string, fn ModuleFn) *Module {
	return &Module{
		Name: name,
		fn:   fn,
	}
}

func (m Module) Execute(ctx goka.Context, message *interface{}, params []string) FilterResult {
	cbCtx := cbContext{GokaCtx: ctx}

	m.fn(&cbCtx, *message, params)

	if cbCtx.Discarded {
		return Discard
	}

	if cbCtx.OverridedMessage != nil {
		message = &cbCtx.OverridedMessage
	}

	return Allow
}
