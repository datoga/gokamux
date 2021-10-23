package transformer

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/builtin/jq/common"
	"github.com/datoga/gokamux/modules/model"
)

func init() {
	modules.MustRegister("jq-transformer", &JQTransformer{})
}

type JQTransformer struct {
	common.JQ
}

func (transformer JQTransformer) Process(ctx model.Context, msg string) {
	v, err := transformer.Eval(ctx.GokaContext().Context(), msg)

	if err != nil {
		ctx.Err(err)
		return
	}

	s, ok := v.(string)

	if !ok {
		ctx.Err(fmt.Errorf("failed decoding JQ result, expecting string, got %T for %v", s, s))
		return
	}

	ctx.OverrideMessage(s)
}
