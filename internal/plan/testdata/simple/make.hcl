default_goal = ["test.txt"]

rule {
  target  = "test.txt"
  command = "echo hello > test.txt"
}
