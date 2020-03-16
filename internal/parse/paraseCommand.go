package parse

import (
	"github.com/pkg/errors"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/hashicorp/hcl/v2"
)

var (
	commandSchema = hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			hcl.AttributeSchema{
				Name:     "command",
				Required: false,
			},
			hcl.AttributeSchema{
				Name:     "dependencies",
				Required: false,
			},
			hcl.AttributeSchema{
				Name:     "environment",
				Required: false,
			},
		},
	}
)

func constructCommand(blk *hcl.Block, ctx *hcl.EvalContext) (*definition.Command, error) {
	con, diag := blk.Body.Content(&commandSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	var c definition.Command
	c.Name = blk.Labels[0]

	for name, attr := range con.Attributes {
		switch name {
		case "environment":
			environment, err := evaluateStringMap(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate environment")
				return nil, err
			}

			c.Environment = environment
		case "command":
			command, err := evaluateString(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate command")
				return nil, err
			}

			c.Command = command
		case "dependencies":
			dependencies, err := evaluateStringArray(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate dependencies")
				return nil, err
			}

			c.Dependencies = dependencies
		}
	}

	return &c, nil
}
