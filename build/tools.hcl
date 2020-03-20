var {
  toolPaths = [
    "github.com/hashicorp/hcl/v2/cmd/hclfmt",
  ]
  tools = {
    for toolPath in var.toolPaths : basename(toolPath) => {
      path: toolPath,
      binPath: path("bin", basename(toolPath)),
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
