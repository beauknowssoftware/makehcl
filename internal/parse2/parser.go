package parse2

import (
	"fmt"

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

	return f, diag
}

func (p *Parser) enumerateImportBlocks(ctx *hcl.EvalContext) hcl.Diagnostics {
	var result hcl.Diagnostics

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

		diag = f.enumerateImportBlocks(ctx)
		if diag.HasErrors() {
			result = result.Extend(diag)
		}

		newFilenames := make([]string, 0, len(f.ImportBlocks))

		for _, imp := range f.ImportBlocks {
			if imp.File == nil {
				continue
			}

			newFilenames = append(newFilenames, imp.File.Value)
		}

		unprocessedFiles = append(unprocessedFiles, newFilenames...)
	}

	return result
}

type importCycleDetector struct {
	visited  map[string]bool
	visiting map[string]bool
	files    map[string]*File
}

func (d importCycleDetector) findImportCycles(filename string) hcl.Diagnostics {
	if d.visited[filename] {
		return nil
	}

	var result hcl.Diagnostics

	d.visiting[filename] = true

	f := d.files[filename]
	for _, imp := range f.ImportBlocks {
		if imp.File == nil {
			continue
		}

		if d.visiting[imp.File.Value] {
			diag := hcl.Diagnostic{
				Summary:     "Import cycle detected",
				Detail:      fmt.Sprintf("Cycle occurred when importing %v", imp.File.Value),
				Severity:    hcl.DiagError,
				Subject:     &imp.File.attribute.Range,
				Expression:  imp.File.attribute.Expr,
				EvalContext: imp.File.ctx,
			}

			result = result.Append(&diag)
		} else {
			diag := d.findImportCycles(imp.File.Value)
			result = result.Extend(diag)
		}
	}

	d.visited[filename] = true
	d.visiting[filename] = false

	return result
}

func (p Parser) findImportCycles() (result hcl.Diagnostics) {
	d := importCycleDetector{
		visited:  make(map[string]bool),
		visiting: make(map[string]bool),
		files:    p.Files,
	}

	return d.findImportCycles(p.Options.Filename)
}

func (p *Parser) Parse() hcl.Diagnostics {
	p.init()

	if diag := p.enumerateImportBlocks(nil); diag.HasErrors() {
		return diag
	}

	if p.Options.StopAfterStage == StopAfterImports {
		return nil
	}

	if diag := p.findImportCycles(); diag.HasErrors() {
		return diag
	}

	return nil
}
