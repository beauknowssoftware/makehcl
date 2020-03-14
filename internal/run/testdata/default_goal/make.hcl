default_goal = [
  "test.txt",
  "test2.txt",
]

rule {
  target = "test.txt"
  command = "touch test.txt"
}

rule {
  target = "test2.txt"
  command = "touch test2.txt"
}

rule {
  target = "test3.txt"
  command = "touch test3.txt"
}
