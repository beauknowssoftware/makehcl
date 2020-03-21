package parse2

import "github.com/hashicorp/hcl/v2"

type varBlock struct {
	block      *hcl.Block
	attributes []attribute
}

func (blk *varBlock) initAttributes(gs scope) hcl.Diagnostics {
	attrs, result := blk.block.Body.JustAttributes()

	blk.attributes = make([]attribute, 0, len(attrs))

	for name, attr := range attrs {
		attr := attribute{
			set:      setOnObject("var", name),
			fillable: &generic{attribute: attr},
			scope:    gs,
		}
		blk.attributes = append(blk.attributes, attr)
	}

	return result
}
