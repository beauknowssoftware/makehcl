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
)

type File struct {
	Name            string
	hclFile         *hcl.File
	unprocessedBody hcl.Body
	content         *hcl.BodyContent
	ImportBlocks    []*ImportBlock
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
