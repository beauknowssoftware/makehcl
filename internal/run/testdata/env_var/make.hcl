var {
  target = env.TARGET
}

rule {
  target  = var.target
  command = "echo hello > test.txt"
}
