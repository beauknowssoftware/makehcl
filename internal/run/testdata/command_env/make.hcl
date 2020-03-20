command touch {
  command = "echo hello > $TARGET"
  environment = {
    TARGET : "test.txt"
  }
}
