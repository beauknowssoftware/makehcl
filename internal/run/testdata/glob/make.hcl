default_goal = ["1.txt", "2.txt"]

dynamic rule {
  for_each = glob("*.original")
  target = join("", [filename(rule), ".txt"])
  command = "cp ${rule} ${target}"
}
