default_goal = rule.tests

dynamic rule {
  name = "tests"
  for_each = ["test", "test2"]
  target = join(".", [rule, "txt"])
  command = "touch ${target}"
}
