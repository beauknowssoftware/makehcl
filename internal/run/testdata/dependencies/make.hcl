rule {
  target = "test.txt"
  command = "touch test.txt"
  dependencies = [
    "test2.txt"
  ]
}

rule {
  target = "test2.txt"
  command = "touch test2.txt"
}
