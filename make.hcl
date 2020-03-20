default_goal = concat(rule.bins, rule.env_bins, ".lint")

// show all commands if in debug mode
var {
  is_debug = exists(env, "DEBUG")
}
opts {
  shell       = "/bin/bash"
  shell_flags = var.is_debug ? "-xuec" : "-uec"
}

// common go build variables
env {
  GOSUMDB = "off"
  GOPROXY = "direct"
}

// go prebuild
rule {
  target       = ".import"
  tee_target   = true
  command      = "goimports -w ."
  dependencies = concat(glob("**.go"), "go.mod", "go.sum")
}
rule {
  target       = ".test"
  tee_target   = true
  dependencies = concat(".import", glob("**/testdata/**"))
  command      = "go test -count=1 ./..."
}
rule {
  target     = ".lint"
  tee_target = true
  dependencies = ".test"
  command = "golangci-lint run --fix"
}

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

  target = cmd.bin
  dependencies = ".test"
  command = "go build -o ${target} ./${cmd.path}"
}

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

  target = cmd.bin
  dependencies = ".test"
  command = "go build -o ${target} ./${cmd.path}"

  environment = {
    GOOS   = cmd.env.goos
    GOARCH = cmd.env.goarch
  }
}

dynamic command {
  alias    = "install"
  for_each = var.cmds
  as       = "cmd"

  name = "install_${cmd.name}"
  dependencies = ".test"
  command = "go install ./${cmd.path}"
}

command clean {
  command = "git clean -f -fdX"
}
