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
				Name:     "alias",
				Required: false,
			},
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

type dynamicRule struct {
	alias string
	rules []*definition.Rule
}

func constructDynamicRules(blk *hcl.Block, ctx *hcl.EvalContext) (*dynamicRule, error) {
	con, body, diag := blk.Body.PartialContent(&dynamicRuleListSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	forEach, err := evaluateValueArray(con.Attributes["for_each"].Expr, ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to evaluate for_each")
		return nil, err
	}

	var dr dynamicRule

	as := "rule"
	if asVal, hasAs := con.Attributes["as"]; hasAs {
		as, err = evaluateString(asVal.Expr, ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to evaluate as")
			return nil, err
		}
	}

	if aliasVal, hasName := con.Attributes["alias"]; hasName {
		dr.alias, err = evaluateString(aliasVal.Expr, ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to evaluate alias")
			return nil, err
		}
	}

	dr.rules = make([]*definition.Rule, 0, len(forEach))

	for _, each := range forEach {
		ectx := ctx.NewChild()
		ectx.Variables = map[string]cty.Value{
			as: each,
		}

		r, err := fillRule(body, ectx)
		if err != nil {
			return nil, err
		}

		dr.rules = append(dr.rules, r)
	}

	return &dr, nil
}
