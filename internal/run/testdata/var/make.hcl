var {
  target = "test.txt"
}

rule {
  target = var.target
  command = "echo hello > test.txt"
}
