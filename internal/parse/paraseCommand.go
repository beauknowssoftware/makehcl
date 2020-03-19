package parse

import (
	"github.com/pkg/errors"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/hashicorp/hcl/v2"
)

const (
	commandAttributeName      = "command"
	dependenciesAttributeName = "dependencies"
	environmentAttributeName  = "environment"
)

var (
	commandSchema = hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     commandAttributeName,
				Required: false,
			},
			{
				Name:     dependenciesAttributeName,
				Required: false,
			},
			{
				Name:     environmentAttributeName,
				Required: false,
			},
		},
	}
)

func fillCommand(name string, body hcl.Body, ctx *hcl.EvalContext) (*definition.Command, error) {
	con, diag := body.Content(&commandSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	var c definition.Command
	c.Name = name

	for name, attr := range con.Attributes {
		switch name {
		case environmentAttributeName:
			environment, err := evaluateStringMap(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrapf(err, "failed to evaluate %v", environmentAttributeName)
				return nil, err
			}

			c.Environment = environment
		case commandAttributeName:
			command, err := evaluateString(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrapf(err, "failed to evaluate %v", commandAttributeName)
				return nil, err
			}

			c.Command = command
		case dependenciesAttributeName:
			dependencies, err := evaluateStringArray(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrapf(err, "failed to evaluate %v", dependenciesAttributeName)
				return nil, err
			}

			c.Dependencies = dependencies
		}
	}

	return &c, nil
}

func constructCommand(blk *hcl.Block, ctx *hcl.EvalContext) (*definition.Command, error) {
	return fillCommand(blk.Labels[0], blk.Body, ctx)
}
