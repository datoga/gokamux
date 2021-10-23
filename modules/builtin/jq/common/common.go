package common

import (
	"fmt"

	"github.com/itchyny/gojq"
	"golang.org/x/net/context"
)

type JQ struct {
	query gojq.Query
}

func (jq *JQ) Init(params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("expression required")
	}

	query, err := gojq.Parse(params[0])

	if err != nil {
		return err
	}

	jq.query = *query

	return nil
}

func (jq JQ) Eval(ctx context.Context, msg string) (interface{}, error) {
	iter := jq.query.RunWithContext(ctx, msg)

	v, ok := iter.Next()

	if !ok {
		return v, fmt.Errorf("no results for JQ")
	}

	if err, ok := v.(error); ok {
		return v, fmt.Errorf("failed decoding JQ result with error %v", err)
	}

	return v, nil
}
