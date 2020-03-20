rule {
  target       = "test.txt"
  command      = "touch ${target}"
  dependencies = "test2.txt"
}
