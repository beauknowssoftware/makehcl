package parse2

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type varBlock struct {
	block      *hcl.Block
	attributes []attribute
}

func (blk *varBlock) initAttributes(gs scope) hcl.Diagnostics {
	attrs, result := blk.block.Body.JustAttributes()

	blk.attributes = make([]attribute, 0, len(attrs))

	for name, attr := range attrs {
		gen := generic{attribute: attr}
		attr := attribute{
			name:         fmt.Sprintf("var.%v", name),
			set:          setOnObject("var", name),
			fillable:     &gen,
			scope:        gs,
			dependencies: getDependencies("", gen.attribute.Expr),
		}
		blk.attributes = append(blk.attributes, attr)
	}

	return result
}
