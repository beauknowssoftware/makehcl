default_goal = ["test.txt.1", "test.txt.2", "test.txt.3", "test.txt.4"]

var {
  target1 = join(".", [var.target, "1"])
}

var {
  target2 = join(".", [var.target, "2"])
  target  = "test.txt"
  target3 = join(".", [var.target, "3"])
}

var {
  target4 = join(".", [var.target, "4"])
}

rule {
  target  = var.target1
  command = "echo hello1 > ${target}"
}

rule {
  target  = var.target2
  command = "echo hello2 > ${target}"
}

rule {
  target  = var.target3
  command = "echo hello3 > ${target}"
}

rule {
  target  = var.target4
  command = "echo hello4 > ${target}"
}
