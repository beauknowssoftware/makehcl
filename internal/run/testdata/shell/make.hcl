var {
  target = shell("echo test.txt")
}

rule {
  target  = var.target
  command = "touch ${target}"
}
