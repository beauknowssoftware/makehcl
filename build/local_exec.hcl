// local executable binaries
var {
  cmds = [for cmd in glob("cmd/*") : {
    path : cmd,
    bin : path("bin/", basename(cmd)),
    name : basename(cmd)
  }]
}
dynamic rule {
  alias    = "bins"
  for_each = var.cmds
  as       = "cmd"

  target       = cmd.bin
  dependencies = ".test"
  command      = "go build -o ${target} ./${cmd.path}"
}
