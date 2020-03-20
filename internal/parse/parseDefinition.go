package parse

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"

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

func fillRuleFromRuleBlock(blk *hcl.Block, d *definition.Definition, ctx *hcl.EvalContext) error {
	r, err := constructRule(blk, ctx)
	if err != nil {
		return err
	}

	d.AddRule(r)

	return nil
}

type dynamicTarget struct {
	targetType    string
	alias         string
	targets       []cty.Value
	targetStrings []string
}

func (dt *dynamicTarget) addTarget(t string) {
	dt.targetStrings = append(dt.targetStrings, t)
	dt.targets = append(dt.targets, cty.StringVal(t))
}

func fillFromDynamicBlock(blk *hcl.Block, d *definition.Definition, ctx *hcl.EvalContext) (*dynamicTarget, error) {
	switch blk.Labels[0] {
	case ruleBlockType:
		dy, err := constructDynamicRules(blk, ctx)
		if err != nil {
			return nil, err
		}

		var dt dynamicTarget
		dt.alias = dy.alias
		dt.targetType = ruleBlockType
		dt.targets = make([]cty.Value, 0, len(dy.rules))

		for _, dr := range dy.rules {
			d.AddRule(dr)
			dt.addTarget(dr.Target)
		}

		return &dt, nil
	case commandBlockType:
		dy, err := constructDynamicCommands(blk, ctx)
		if err != nil {
			return nil, err
		}

		var dt dynamicTarget
		dt.alias = dy.alias
		dt.targetType = commandBlockType
		dt.targets = make([]cty.Value, 0, len(dy.commands))

		for _, dc := range dy.commands {
			d.AddCommand(dc)
			dt.addTarget(dc.Name)
		}

		return &dt, nil
	default:
		return nil, fmt.Errorf("unknown dynamic type %v", blk.Labels[0])
	}
}

func fillRuleFromCommandBlock(blk *hcl.Block, d *definition.Definition, ctx *hcl.EvalContext) error {
	c, err := constructCommand(blk, ctx)
	if err != nil {
		return err
	}

	d.AddCommand(c)

	return nil
}

func fillRules(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	for _, blk := range con.Blocks {
		switch blk.Type {
		case ruleBlockType:
			if err := fillRuleFromRuleBlock(blk, d, ctx); err != nil {
				return err
			}
		case commandBlockType:
			if err := fillRuleFromCommandBlock(blk, d, ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func fillDynamicRules(con *hcl.BodyContent, d *definition.Definition, ctx *hcl.EvalContext) error {
	rules := make(map[string]cty.Value)
	commands := make(map[string]cty.Value)

	for _, blk := range con.Blocks {
		if blk.Type == dynamicBlockType {
			dt, err := fillFromDynamicBlock(blk, d, ctx)
			if err != nil {
				return err
			}

			if dt.alias != "" && dt.targetType == commandBlockType {
				commands[dt.alias] = cty.ListVal(dt.targets)
			}

			if dt.alias != "" && dt.targetType == ruleBlockType {
				rules[dt.alias] = cty.ListVal(dt.targets)
			}

			if dt.alias != "" {
				d.AddCommand(&definition.Command{
					Name:         dt.alias,
					Dependencies: dt.targetStrings,
				})
			}
		}
	}

	if len(rules) > 0 {
		ctx.Variables["rule"] = cty.MapVal(rules)
	}

	if len(commands) > 0 {
		ctx.Variables["command"] = cty.MapVal(commands)
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

	if err := fillDynamicRules(con, &d, ctx); err != nil {
		return nil, err
	}

	if err := fillRules(con, &d, ctx); err != nil {
		return nil, err
	}

	if err := fillDefaultGoal(con, &d, ctx); err != nil {
		return nil, err
	}

	return &d, nil
}
