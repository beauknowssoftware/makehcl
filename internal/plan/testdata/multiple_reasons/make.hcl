default_goal = ["test.txt"]

rule {
  target       = "test.txt"
  command      = "touch ${target}"
  dependencies = ["test2.txt", "test3.txt"]
}

rule {
  target  = "test2.txt"
  command = "touch ${target}"
}

rule {
  target  = "test3.txt"
  command = "touch ${target}"
}