package parse

import (
	"github.com/pkg/errors"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

var (
	dynamicCommandListSchema = hcl.BodySchema{
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
			{
				Name:     "name",
				Required: true,
			},
		},
	}
)

type dynamicCommand struct {
	alias    string
	commands []*definition.Command
}

func constructDynamicCommands(blk *hcl.Block, ctx *hcl.EvalContext) (*dynamicCommand, error) {
	con, body, diag := blk.Body.PartialContent(&dynamicCommandListSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	forEach, err := evaluateValueArray(con.Attributes["for_each"].Expr, ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to evaluate for_each")
		return nil, err
	}

	var dr dynamicCommand

	as := "command"
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

	dr.commands = make([]*definition.Command, 0, len(forEach))

	for _, each := range forEach {
		ectx := ctx.NewChild()
		ectx.Variables = map[string]cty.Value{
			as: each,
		}

		name, err := evaluateString(con.Attributes["name"].Expr, ectx)
		if err != nil {
			err = errors.Wrap(err, "failed to evaluate name")
			return nil, err
		}

		ectx.Variables["name"] = cty.StringVal(name)

		r, err := fillCommand(name, body, ectx)
		if err != nil {
			return nil, err
		}

		dr.commands = append(dr.commands, r)
	}

	return &dr, nil
}
