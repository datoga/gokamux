package model

import "github.com/lovoo/goka"

type Context interface {
	Discard()
	OverrideMessage(msg string)
	Err(error)
	GokaContext() goka.Context
}
