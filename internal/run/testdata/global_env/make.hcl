env {
  TARGET = "test.txt"
}

rule {
  target = "test.txt"
  command = "echo hello > $TARGET"
}
