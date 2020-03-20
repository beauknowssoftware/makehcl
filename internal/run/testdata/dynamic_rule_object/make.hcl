default_goal = ["test.txt", "test2.txt"]

var {
  targets = {
    "test" : "test.txt",
    "test2" : "test2.txt",
  }
}

dynamic rule {
  for_each = var.targets
  target   = rule
  command  = "touch ${rule}"
}
