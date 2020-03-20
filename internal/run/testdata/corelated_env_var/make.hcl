var {
  target  = join(".", [env.NAME, env.EXT])
  target2 = env.TARGET
}

env {
  NAME   = join("_", [env.NAME, "new"])
  TARGET = var.target
  VALUE  = join("_", [env.VALUE, "new"])
}

rule {
  target  = var.target2
  command = "echo $VALUE > ${target}"
}
