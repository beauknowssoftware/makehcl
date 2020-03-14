opts {
  shell = "touch"
  shell_flags = ""
}

rule {
  target = "test.txt"
  command = "${target}"
}
