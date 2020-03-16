default_goal = concat(var.bins, var.env_bins)

var {
  is_debug = exists(env, "DEBUG")
  go_deps = concat(
    glob("**.go"),
    ["go.mod", "go.sum", "test", "lint"]
  )
}

opts {
  shell = "/bin/bash"
  shell_flags = var.is_debug ? "-xuec" : "-uec"
}

env {
  GOSUMDB = "off"
  GOPROXY = "direct"
}

var {
  cmds = [for cmd in glob("cmd/*") : { path: cmd, bin: path("bin/", basename(cmd)) }]
  bins = [for cmd in var.cmds : cmd.bin]
}

dynamic rule {
  for_each = var.cmds
  as = "cmd"

  target = cmd.bin
  dependencies = var.go_deps
  command = "go build -o ${target} ./${cmd.path}"
}

var {
  go_envs = [
    { goos: "darwin", goarch: "386" },
    { goos: "darwin", goarch: "amd64" },
    { goos: "linux", goarch: "386" },
    { goos: "linux", goarch: "amd64" },
  ]
  env_cmds = flatten([
    for cmd in var.cmds : [
      for env in var.go_envs : { path: cmd.path, bin: path("bin/", env.goos, env.goarch, basename(cmd.path)), env: env, }
    ]
  ])
  env_bins = [for cmd in var.env_cmds : cmd.bin]
}

dynamic rule {
  for_each = var.env_cmds
  as = "cmd"

  target = cmd.bin
  dependencies = var.go_deps
  command = "go build -o ${target} ./${cmd.path}"

  environment = {
    GOOS = cmd.env.goos
    GOARCH = cmd.env.goarch
  }
}

command lint { command = "golangci-lint run --fix" }
command test {
  dependencies = ["import"]
  command = "go test ./..."
}
command import { command = "goimports -w ." }
command clean { command = "git clean -f -fdX" }
