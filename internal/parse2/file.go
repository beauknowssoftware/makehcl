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
	ruleBlockHeaderSchema = hcl.BlockHeaderSchema{Type: "rule"}
	ruleStageSchema       = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			ruleBlockHeaderSchema,
		},
	}
)

type File struct {
	Name            string
	hclFile         *hcl.File
	unprocessedBody hcl.Body
	content         *hcl.BodyContent
	ImportBlocks    []*ImportBlock
	RuleBlocks      []*RuleBlock
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

func (f *File) enumerateAttributes(gs scope) (result hcl.Diagnostics) {
	if f.unprocessedBody == nil {
		return
	}

	f.content, f.unprocessedBody, result = f.unprocessedBody.PartialContent(ruleStageSchema)

	if f.content == nil {
		return
	}

	for _, blk := range f.content.Blocks {
		if blk.Type != ruleBlockHeaderSchema.Type {
			continue
		}

		rBlk := RuleBlock{block: blk}
		f.RuleBlocks = append(f.RuleBlocks, &rBlk)

		if diag := rBlk.initAttributes(gs); diag.HasErrors() {
			result = result.Extend(diag)
		}

		f.attributes = append(f.attributes, rBlk.attributes...)
	}

	return
}
