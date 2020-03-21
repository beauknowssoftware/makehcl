package parse2

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

const (
	defaultFilename = "make.hcl"
)

type Parser struct {
	Options   Options
	hclParser *hclparse.Parser
	Definition
}

func (p *Parser) init() {
	if p.Options.Filename == "" {
		p.Options.Filename = defaultFilename
	}

	p.hclParser = hclparse.NewParser()
}

func (p *Parser) readFile(filename string) (*File, hcl.Diagnostics) {
	hf, diag := p.hclParser.ParseHCLFile(filename)

	f := p.addFile(filename, hf)

	if diag.HasErrors() {
		return f, diag
	}

	diag = f.enumerateContents()

	return f, diag
}

func (p *Parser) getBlocks(ctx *hcl.EvalContext) (result hcl.Diagnostics) {
	unprocessedFiles := []string{p.Options.Filename}

	for len(unprocessedFiles) > 0 {
		filename := unprocessedFiles[0]
		unprocessedFiles = unprocessedFiles[1:]

		if _, hasReadFile := p.Files[filename]; hasReadFile {
			continue
		}

		f, diag := p.readFile(filename)
		if diag.HasErrors() {
			result = result.Extend(diag)
		}

		newFilenames, diag := f.getImportFilenames(ctx)
		if diag.HasErrors() {
			result = result.Extend(diag)
		}

		unprocessedFiles = append(unprocessedFiles, newFilenames...)
	}

	return
}

func (p *Parser) Parse() hcl.Diagnostics {
	p.init()

	if diag := p.getBlocks(nil); diag.HasErrors() {
		return diag
	}

	return nil
}
