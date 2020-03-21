var {
  target = "test.txt"
}

rule {
  target = var.target
  command = "touch ${target}"
}