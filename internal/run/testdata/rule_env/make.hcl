rule {
  target  = "test.txt"
  command = "echo hello > $TARGET"
  environment = {
    TARGET : "test.txt"
  }
}
