default_goal = [
  "test.txt",
  "test2.txt",
]

rule {
  target = "test.txt"
  command = "echo hello >> test.txt"
  dependencies = ["test3.txt"]
}

rule {
  target = "test2.txt"
  command = "echo hello2 >> test2.txt"
  dependencies = ["test3.txt"]
}

rule {
  target = "test3.txt"
  command = "echo hello3 >> test3.txt"
}
