default_goal = [var.secondTarget ? "test2.txt" : "test.txt"]

var {
  secondTarget = exists(env, "SECOND_TARGET")
}

rule {
  target  = "test.txt"
  command = "touch ${target}"
}

rule {
  target  = "test2.txt"
  command = "touch ${target}"
}
