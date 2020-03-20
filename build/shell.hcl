// show all commands if in debug mode
var {
  is_debug = exists(env, "DEBUG")
}
opts {
  shell       = "/bin/bash"
  shell_flags = var.is_debug ? "-xuec" : "-uec"
}
