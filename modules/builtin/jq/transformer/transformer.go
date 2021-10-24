package transformer

import (
	"fmt"

	"github.com/datoga/gokamux/modules"
	"github.com/datoga/gokamux/modules/builtin/jq/common"
	"github.com/datoga/gokamux/modules/model"
)

func init() {
	modules.MustRegister("jq-transformer", &jqTransformer{})
}

type jqTransformer struct {
	common.JQ
}

func (transformer jqTransformer) Process(ctx model.Context, msg string) error {
	v, err := transformer.Eval(ctx.GokaContext().Context(), msg)

	if err != nil {
		return err
	}

	s, ok := v.(string)

	if !ok {
		return fmt.Errorf("failed decoding JQ result, expecting string, got %T for %v", s, s)
	}

	ctx.OverrideMessage(s)

	return nil
}
