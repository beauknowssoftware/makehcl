var {
  // toolPaths = split(" ", shell("go list -f '{{ .Imports }}' -tags tools ./internal")
  toolPaths = [
    "github.com/hashicorp/hcl/v2/cmd/hclfmt",
    "golang.org/x/tools/cmd/goimports",
    "github.com/golangci/golangci-lint/cmd/golangci-lint"
  ]
  tools = {
    for toolPath in var.toolPaths : basename(toolPath) => {
      path: toolPath,
      binPath: path("tools", basename(toolPath)),
      name: basename(toolPath),
    }
  }
}

dynamic rule {
  for_each = var.tools
  as = "tool"

  target = tool.binPath
  command = "go build -o ${target} ${tool.path}"
}
