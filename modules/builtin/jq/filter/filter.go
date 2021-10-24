package filter

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/builtin/jq/common"
	"github.com/datoga/gokamux/modules/model"
)

func init() {
	modules.MustRegister("jq-filter", &jqFilter{})
}

type jqFilter struct {
	common.JQ
}

func (filter jqFilter) Process(ctx model.Context, msg string) error {
	v, err := filter.Eval(ctx.GokaContext().Context(), msg)

	if err != nil {
		return err
	}

	b, ok := v.(bool)

	if !ok {
		return fmt.Errorf("failed decoding JQ result, expecting boolean, got %T for %v", b, b)
	}

	if b {
		ctx.Discard()
	}

	return nil
}
