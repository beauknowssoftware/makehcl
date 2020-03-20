dynamic command {
  alias    = "install"
  for_each = var.cmds
  as       = "cmd"

  name         = "install_${cmd.name}"
  dependencies = ".test"
  command      = "go install ./${cmd.path}"
}
