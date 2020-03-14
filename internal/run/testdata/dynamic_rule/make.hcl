default_goal = ["test.txt", "test2.txt"]

dynamic rule {
  for_each = ["test.txt", "test2.txt"]
  target = rule
  command = "touch ${rule}"
}
