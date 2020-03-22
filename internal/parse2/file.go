package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

var (
	importBlockHeaderSchema = hcl.BlockHeaderSchema{Type: "import"}
	importStageSchema       = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			importBlockHeaderSchema,
		},
	}
	ruleBlockHeaderSchema    = hcl.BlockHeaderSchema{Type: "rule"}
	commandBlockHeaderSchema = hcl.BlockHeaderSchema{
		Type:       "command",
		LabelNames: []string{"name"},
	}
	varBlockHeaderSchema       = hcl.BlockHeaderSchema{Type: "var"}
	defaultGoalAttributeSchema = hcl.AttributeSchema{Name: "default_goal"}
	attributeStageSchema       = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			defaultGoalAttributeSchema,
		},
		Blocks: []hcl.BlockHeaderSchema{
			commandBlockHeaderSchema,
			ruleBlockHeaderSchema,
			varBlockHeaderSchema,
		},
	}
)

type File struct {
	Name            string
	DefaultGoal     *StringArray
	hclFile         *hcl.File
	unprocessedBody hcl.Body
	content         *hcl.BodyContent
	scope           scope
	ImportBlocks    []*ImportBlock
	RuleBlocks      []*RuleBlock
	CommandBlocks   []*CommandBlock
	varBlocks       []*varBlock
	attributes      []attribute
}

func (f File) HasContents() bool {
	return f.hclFile != nil
}

func (f *File) enumerateImportBlocks(ctx *hcl.EvalContext) (result hcl.Diagnostics) {
	if f.unprocessedBody == nil {
		return
	}

	f.content, f.unprocessedBody, result = f.unprocessedBody.PartialContent(importStageSchema)

	if f.content == nil {
		return
	}

	for _, blk := range f.content.Blocks {
		if blk.Type != importBlockHeaderSchema.Type {
			continue
		}

		iBlk := ImportBlock{block: blk}
		f.ImportBlocks = append(f.ImportBlocks, &iBlk)

		if diag := iBlk.initAttributes(ctx); diag.HasErrors() {
			result = result.Extend(diag)
		}
	}

	return
}

func (f *File) enumerateAttributes(gs scope) hcl.Diagnostics {
	var result hcl.Diagnostics

	if f.unprocessedBody == nil {
		return result
	}

	f.content, f.unprocessedBody, result = f.unprocessedBody.PartialContent(attributeStageSchema)

	if f.content == nil {
		return result
	}

	f.scope = gs

	for _, attr := range f.content.Attributes {
		if attr.Name == defaultGoalAttributeSchema.Name {
			f.DefaultGoal = &StringArray{attribute: attr}
			attr := attribute{
				name:         "default_goal",
				fillable:     f.DefaultGoal,
				scope:        f.scope,
				dependencies: getDependencies("", f.DefaultGoal.attribute.Expr),
			}
			f.attributes = append(f.attributes, attr)
		}
	}

	for _, blk := range f.content.Blocks {
		switch blk.Type {
		case commandBlockHeaderSchema.Type:
			cBlk := CommandBlock{block: blk}
			f.CommandBlocks = append(f.CommandBlocks, &cBlk)

			if diag := cBlk.initAttributes(gs); diag.HasErrors() {
				result = result.Extend(diag)
			}

			f.attributes = append(f.attributes, cBlk.attributes...)
		case ruleBlockHeaderSchema.Type:
			rBlk := RuleBlock{block: blk}
			f.RuleBlocks = append(f.RuleBlocks, &rBlk)

			if diag := rBlk.initAttributes(gs); diag.HasErrors() {
				result = result.Extend(diag)
			}

			f.attributes = append(f.attributes, rBlk.attributes...)
		case varBlockHeaderSchema.Type:
			vBlk := varBlock{block: blk}
			f.varBlocks = append(f.varBlocks, &vBlk)

			if diag := vBlk.initAttributes(gs); diag.HasErrors() {
				result = result.Extend(diag)
			}

			f.attributes = append(f.attributes, vBlk.attributes...)
		}
	}

	return result
}
