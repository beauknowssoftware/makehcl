default_goal = "test.txt"

rule {
  target = "test.txt"
  command = "touch ${target}"
}