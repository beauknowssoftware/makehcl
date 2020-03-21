package parse2

import "github.com/hashicorp/hcl/v2"

type RuleBlock struct {
	block      *hcl.Block
	content    *hcl.BodyContent
	Target     *String
	Command    *StringArray
	attributes []attribute
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

	blk.attributes = make([]attribute, 0, len(con.Attributes))

	for _, attr := range con.Attributes {
		switch attr.Name {
		case ruleTargetAttributeSchema.Name:
			blk.Target = &String{attribute: attr}
			attr := attribute{blk.Target}
			blk.attributes = append(blk.attributes, attr)
		case ruleCommandAttributeSchema.Name:
			blk.Command = &StringArray{attribute: attr}
			attr := attribute{blk.Command}
			blk.attributes = append(blk.attributes, attr)
		}
	}

	return result
}
