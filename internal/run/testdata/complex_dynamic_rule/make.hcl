default_goal = ["test.txt", "test2.txt"]

var {
  names = ["test", "test2"]
  targets = [
    for name in var.names :
    { target: join("", [name, ".txt"]) }
  ]
}

dynamic rule {
  for_each = var.targets
  target = rule.target
  command = "touch ${rule.target}"
}
