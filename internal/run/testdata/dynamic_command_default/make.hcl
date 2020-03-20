default_goal = command.tests

dynamic command {
  alias = "tests"

  for_each = ["test", "test2"]
  name     = command
  command  = "touch ${name}.txt"
}
