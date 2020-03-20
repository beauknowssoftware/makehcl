default_goal = ["out/1", "out/2"]

dynamic rule {
  for_each = glob("test/*")
  target   = path("out/", basename(rule))
  command  = "mkdir -p $(dirname ${target}); cp ${rule}/original ${target}"
}
