package parse

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/pkg/errors"

	"github.com/beauknowssoftware/makehcl/internal/definition"
)

const (
	commandBlockType = "command"
	ruleBlockType    = "rule"
	dynamicBlockType = "dynamic"
)

var (
	definitionSchema = hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "default_goal",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "opts",
			},
			{
				Type: "env",
			},
			{
				Type: "var",
			},
			{
				Type:       commandBlockType,
				LabelNames: []string{"name"},
			},
			{
				Type: ruleBlockType,
			},
			{
				Type:       dynamicBlockType,
				LabelNames: []string{"type"},
			},
		},
	}
)

func getAllAttributes(blockType string, con *hcl.BodyContent) (map[string]*hcl.Attribute, error) {
	attrs := make(map[string]*hcl.Attribute)

	for _, blk := range con.Blocks {
		if blk.Type == blockType {
			attr, diag := blk.Body.JustAttributes()
			if diag.HasErrors() {
				return nil, errors.Wrapf(diag, "failed to get %v attributes", blockType)
			}

			for k, v := range attr {
				attrs[k] = v
			}
		}
	}

	return attrs, nil
}

func fillGlobals(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	varAttrs, err := getAllAttributes("var", con)
	if err != nil {
		return err
	}

	envAttrs, err := getAllAttributes("env", con)
	if err != nil {
		return err
	}

	envs, err := getGlobals(map[string]map[string]*hcl.Attribute{
		"var": varAttrs,
		"env": envAttrs,
	}, ctx)
	if err != nil {
		return err
	}

	d.GlobalEnvironment = envs

	return nil
}

func fillOpts(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	optsAttrs, err := getAllAttributes("opts", con)
	if err != nil {
		return err
	}

	for name, attr := range optsAttrs {
		switch name {
		case "shell":
			v, err := evaluateString(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate shell opt")
				return err
			}

			d.Shell = v
		case "shell_flags":
			v, err := evaluateString(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate shell_flag opt")
				return err
			}

			d.ShellFlags = &v
		}
	}

	return nil
}

func fillDefaultGoal(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	for name, attr := range con.Attributes {
		if name == "default_goal" {
			defaultGoal, err := evaluateStringArray(attr.Expr, ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to evaluate default_goal")
				return err
			}

			d.SetDefaultGoal(defaultGoal)
		}
	}

	return nil
}

func fillRules(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	for _, blk := range con.Blocks {
		switch blk.Type {
		case ruleBlockType:
			r, err := constructRule(blk, ctx)
			if err != nil {
				return err
			}

			d.AddRule(r)
		case dynamicBlockType:
			switch blk.Labels[0] {
			case "rule":
				dy, err := constructDynamicRules(blk, ctx)
				if err != nil {
					return err
				}

				for _, dy := range dy {
					d.AddRule(dy)
				}
			default:
				return fmt.Errorf("unknown dynamic type %v", blk.Labels[0])
			}
		case commandBlockType:
			c, err := constructCommand(blk, ctx)
			if err != nil {
				return err
			}

			d.AddCommand(c)
		}
	}

	return nil
}

func constructDefinition(f *hcl.File, ctx *hcl.EvalContext) (*definition.Definition, error) {
	con, diag := f.Body.Content(&definitionSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	var d definition.Definition

	if err := fillGlobals(con, &d, ctx); err != nil {
		return nil, err
	}

	if err := fillOpts(con, &d, ctx); err != nil {
		return nil, err
	}

	if err := fillDefaultGoal(con, &d, ctx); err != nil {
		return nil, err
	}

	if err := fillRules(con, &d, ctx); err != nil {
		return nil, err
	}

	return &d, nil
}
