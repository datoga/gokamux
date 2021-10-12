package gokamux

import "github.com/lovoo/goka"

type Context interface {
	Discard()
	OverrideMessage(interface{})
	Err(error)
	GokaContext() goka.Context
}

type cbContext struct {
	GokaCtx          goka.Context
	Discarded        bool
	OverridedMessage interface{}
	Error            error
}

func (ct *cbContext) Discard() {
	ct.Discarded = true
}

func (ct *cbContext) OverrideMessage(message interface{}) {
	ct.OverridedMessage = message
}

func (ct *cbContext) Err(err error) {
	ct.Error = err
}

func (ct cbContext) GokaContext() goka.Context {
	return ct.GokaCtx
}
