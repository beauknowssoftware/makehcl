default_goal = ["test1.txt", "test2.txt"]

var {
  target1 = "test1.txt"
}

var {
  target2 = "test2.txt"
}

rule {
  target = var.target1
  command = "echo hello1 > ${target}"
}

rule {
  target = var.target2
  command = "echo hello2 > ${target}"
}
