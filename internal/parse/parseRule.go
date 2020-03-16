package parse

import (
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/hashicorp/hcl/v2"
)

var (
	ruleSchema = hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			hcl.AttributeSchema{
				Name:     "environment",
				Required: false,
			},
			hcl.AttributeSchema{
				Name:     "target",
				Required: true,
			},
			hcl.AttributeSchema{
				Name:     "command",
				Required: true,
			},
			hcl.AttributeSchema{
				Name:     "dependencies",
				Required: false,
			},
		},
	}
)

func fillRule(body hcl.Body, ctx *hcl.EvalContext) (*definition.Rule, error) {
	con, diag := body.Content(&ruleSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	target, err := evaluateString(con.Attributes["target"].Expr, ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to evaluate target")
		return nil, err
	}

	ectx := ctx.NewChild()
	ectx.Variables = map[string]cty.Value{
		"target": cty.StringVal(target),
	}

	var r definition.Rule
	r.Target = target

	for name, attr := range con.Attributes {
		switch name {
		case "environment":
			environment, err := evaluateStringMap(attr.Expr, ectx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate environment")
				return nil, err
			}

			r.Environment = environment
		case "command":
			command, err := evaluateString(attr.Expr, ectx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate command")
				return nil, err
			}

			r.Command = command
		case "dependencies":
			dependencies, err := evaluateStringArray(attr.Expr, ectx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate dependencies")
				return nil, err
			}

			r.Dependencies = dependencies
		}
	}

	return &r, nil
}

func constructRule(blk *hcl.Block, ctx *hcl.EvalContext) (*definition.Rule, error) {
	return fillRule(blk.Body, ctx)
}
