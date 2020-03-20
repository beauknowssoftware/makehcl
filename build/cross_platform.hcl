// cross platform binaries
var {
  go_envs = [
    {
      goos : "darwin",
      goarch : "386"
    },
    {
      goos : "darwin",
      goarch : "amd64"
    },
    {
      goos : "linux",
      goarch : "386"
    },
    {
      goos : "linux",
      goarch : "amd64"
    },
  ]
  env_cmds = flatten([
  for cmd in var.cmds : [
  for env in var.go_envs : {
    path : cmd.path,
    bin : path("bin/", env.goos, env.goarch, basename(cmd.path)),
    env : env,
  }
  ]
  ])
}
dynamic rule {
  alias    = "env_bins"
  for_each = var.env_cmds
  as       = "cmd"

  target       = cmd.bin
  dependencies = ".test"
  command      = "go build -o ${target} ./${cmd.path}"

  environment = {
    GOOS   = cmd.env.goos
    GOARCH = cmd.env.goarch
  }
}
