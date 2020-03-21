package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

const (
	importBlockType          = "import"
	commandBlockType         = "command"
	varBlockType             = "var"
	envBlockType             = "env"
	optsBlockType            = "opts"
	ruleBlockType            = "rule"
	dynamicBlockType         = "dynamic"
	defaultGoalAttributeName = "default_goal"
)

var (
	fileSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: defaultGoalAttributeName},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: importBlockType},
			{Type: varBlockType},
			{Type: envBlockType},
			{Type: optsBlockType},
			{Type: ruleBlockType},
			{
				Type:       dynamicBlockType,
				LabelNames: []string{"Type"},
			},
			{
				Type:       commandBlockType,
				LabelNames: []string{"Name"},
			},
		},
	}
)

type File struct {
	Name         string
	hclFile      *hcl.File
	content      *hcl.BodyContent
	ImportBlocks []*ImportBlock
}

func (f *File) enumerateContents() hcl.Diagnostics {
	con, diag := f.hclFile.Body.Content(fileSchema)

	if con == nil {
		return diag
	}

	f.content = con

	for _, blk := range con.Blocks {
		switch blk.Type {
		case commandBlockType:
		case varBlockType:
		case envBlockType:
		case optsBlockType:
		case ruleBlockType:
		case dynamicBlockType:
		case importBlockType:
			iBlk := ImportBlock{block: blk}
			f.ImportBlocks = append(f.ImportBlocks, &iBlk)
		}
	}

	return diag
}

func (f File) getImportFilenames(ctx *hcl.EvalContext) (result []string, resultDiag hcl.Diagnostics) {
	result = make([]string, 0, len(f.ImportBlocks))

	for _, blk := range f.ImportBlocks {
		if diag := blk.initAttributes(ctx); diag.HasErrors() {
			resultDiag = resultDiag.Extend(diag)
			continue
		}

		if blk.File == nil {
			continue
		}

		result = append(result, blk.File.Value)
	}

	return
}
