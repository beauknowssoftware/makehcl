default_goal = ["test"]

dynamic rule {
  alias = "tests"
  for_each = ["test", "test2"]
  target = join(".", [rule, "txt"])
  command = "touch ${target}"
}

command test {
  dependencies = rule.tests
}