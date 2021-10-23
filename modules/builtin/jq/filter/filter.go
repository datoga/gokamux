package filter

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/builtin/jq/common"
	"github.com/datoga/gokamux/modules/model"
)

func init() {
	modules.MustRegister("jq-filter", &JQFilter{})
}

type JQFilter struct {
	common.JQ
}

func (filter JQFilter) Process(ctx model.Context, msg string) {
	v, err := filter.Eval(ctx.GokaContext().Context(), msg)

	if err != nil {
		ctx.Err(err)
		return
	}

	b, ok := v.(bool)

	if !ok {
		ctx.Err(fmt.Errorf("failed decoding JQ result, expecting boolean, got %T for %v", b, b))
		return
	}

	if b {
		ctx.Discard()
	}
}
