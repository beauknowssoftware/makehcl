var {
  toolPaths = split(" ", shell("go list -f '{{ join .Imports \" \" }}' -tags tools ./internal"))
  tools = [
    for toolPath in var.toolPaths : {
      path: toolPath,
      binPath: path("tools", basename(toolPath)),
      name: basename(toolPath),
    }
  ]
}

dynamic rule {
  for_each = var.tools
  as = "tool"

  target = tool.binPath
  command = "go build -o ${target} ${tool.path}"
}
