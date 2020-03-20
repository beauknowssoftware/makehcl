// +build tools

package internal

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/hcl/v2/cmd/hclfmt"
	_ "golang.org/x/tools/cmd/goimports"
)
