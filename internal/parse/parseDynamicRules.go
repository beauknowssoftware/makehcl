package parse

import (
	"github.com/pkg/errors"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

var (
	dynamicRuleListSchema = hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "for_each",
				Required: true,
			},
			{
				Name:     "as",
				Required: false,
			},
		},
	}
)

func constructDynamicRules(blk *hcl.Block, ctx *hcl.EvalContext) ([]*definition.Rule, error) {
	con, body, diag := blk.Body.PartialContent(&dynamicRuleListSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	forEach, err := evaluateValueArray(con.Attributes["for_each"].Expr, ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to evaluate for_each")
		return nil, err
	}

	as := "rule"
	if asVal, hasAs := con.Attributes["as"]; hasAs {
		as, err = evaluateString(asVal.Expr, ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to evaluate as")
			return nil, err
		}
	}

	result := make([]*definition.Rule, 0, len(forEach))

	for _, each := range forEach {
		ectx := ctx.NewChild()
		ectx.Variables = map[string]cty.Value{
			as: each,
		}

		r, err := fillRule(body, ectx)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}
