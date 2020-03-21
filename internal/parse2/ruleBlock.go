package parse2

import "github.com/hashicorp/hcl/v2"

type RuleBlock struct {
	block      *hcl.Block
	content    *hcl.BodyContent
	Target     *String
	Command    *StringArray
	attributes []fillable
}

var (
	ruleTargetAttributeSchema  = hcl.AttributeSchema{Name: "target", Required: true}
	ruleCommandAttributeSchema = hcl.AttributeSchema{Name: "command", Required: true}
	ruleSchema                 = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			ruleTargetAttributeSchema,
			ruleCommandAttributeSchema,
		},
	}
)

func (blk *RuleBlock) initAttributes() hcl.Diagnostics {
	con, result := blk.block.Body.Content(ruleSchema)

	if con == nil {
		return result
	}

	blk.content = con

	blk.attributes = make([]fillable, 0, len(con.Attributes))

	for _, attr := range con.Attributes {
		switch attr.Name {
		case ruleTargetAttributeSchema.Name:
			blk.Target = &String{attribute: attr}
			blk.attributes = append(blk.attributes, blk.Target)
		case ruleCommandAttributeSchema.Name:
			blk.Command = &StringArray{attribute: attr}
			blk.attributes = append(blk.attributes, blk.Command)
		}
	}

	return result
}
