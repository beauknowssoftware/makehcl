rule {
  target = "test.txt"
  command = "cp test2.txt test.txt"
  dependencies = [
    "test2.txt"
  ]
}

