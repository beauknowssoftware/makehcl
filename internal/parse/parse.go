package parse

import (
	"os"
	"strings"

	"github.com/beauknowssoftware/makehcl/internal/definition"
	"github.com/beauknowssoftware/makehcl/internal/functions"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultFilename = "make.hcl"
)

type Options struct {
	Filename string
}

func envValue() cty.Value {
	env := make(map[string]cty.Value)
	for _, e := range os.Environ() {
		p := strings.SplitN(e, "=", 2)
		env[p[0]] = cty.StringVal(p[1])
	}
	return cty.MapVal(env)
}

func createBaseContext() (*hcl.EvalContext, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get working directory")
	}

	var baseContext = hcl.EvalContext{
		Functions: functions.GetFunctions(curdir),
		Variables: make(map[string]cty.Value),
	}

	baseContext.Variables["env"] = envValue()

	return &baseContext, nil
}

func Parse(o Options) (*definition.Definition, error) {
	if o.Filename == "" {
		o.Filename = defaultFilename
	}

	p := hclparse.NewParser()
	f, diag := p.ParseHCLFile(o.Filename)
	if diag.HasErrors() {
		return nil, diag
	}

	ctx, err := createBaseContext()
	if err != nil {
		return nil, err
	}

	return constructDefinition(f, ctx)
}
