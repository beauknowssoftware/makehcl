package parse2

import "github.com/hashicorp/hcl/v2"

type ImportBlock struct {
	block   *hcl.Block
	content *hcl.BodyContent
	File    *StringAttribute
}

const (
	importFileAttributeName = "file"
)

var (
	importSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: importFileAttributeName, Required: true},
		},
	}
)

func (blk *ImportBlock) initAttributes(ctx *hcl.EvalContext) hcl.Diagnostics {
	con, result := blk.block.Body.Content(importSchema)

	if con == nil {
		return result
	}

	blk.content = con

	for _, attr := range con.Attributes {
		if attr.Name == importFileAttributeName {
			if diag := blk.setFile(attr, ctx); diag.HasErrors() {
				result = result.Extend(diag)
			}
		}
	}

	return result
}

func (blk *ImportBlock) setFile(attr *hcl.Attribute, ctx *hcl.EvalContext) (diag hcl.Diagnostics) {
	blk.File, diag = newStringAttribute(attr, ctx)
	return
}
