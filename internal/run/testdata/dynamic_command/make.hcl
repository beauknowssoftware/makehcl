default_goal = ["test", "test2"]

dynamic command {
  for_each = ["test", "test2"]
  name     = command
  command  = "touch ${name}.txt"
}
