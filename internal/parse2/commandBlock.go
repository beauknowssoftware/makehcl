package parse2

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type CommandBlock struct {
	block      *hcl.Block
	content    *hcl.BodyContent
	Name       string
	Command    *StringArray
	attributes []attribute
	scope      scope
}

var (
	commandCommandAttributeSchema = hcl.AttributeSchema{Name: "command", Required: true}
	commandSchema                 = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			commandCommandAttributeSchema,
		},
	}
)

func (blk *CommandBlock) initAttributes(gs scope) hcl.Diagnostics {
	con, result := blk.block.Body.Content(commandSchema)

	if con == nil {
		return result
	}

	blk.content = con

	blk.attributes = make([]attribute, 0, len(con.Attributes))

	blk.scope = &nestedScope{outer: gs}

	local := fmt.Sprintf("command.%v", blk.block.DefRange)

	blk.Name = blk.block.Labels[0]

	for _, attr := range con.Attributes {
		if attr.Name == commandCommandAttributeSchema.Name {
			blk.Command = &StringArray{attribute: attr}
			attr := attribute{
				name:         fmt.Sprintf("%v.command", local),
				fillable:     blk.Command,
				scope:        blk.scope,
				dependencies: getDependencies(local, blk.Command.attribute.Expr),
			}
			blk.attributes = append(blk.attributes, attr)
		}
	}

	return result
}
