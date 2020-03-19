default_goal = ["test", "test2"]

dynamic command {
  for_each = ["test", "test2"]
  as = "name2"

  name = name2
  command = "touch ${name}.txt"
}
