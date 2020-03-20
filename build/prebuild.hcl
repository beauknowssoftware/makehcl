// go prebuild
rule {
  target       = ".import"
  tee_target   = true
  command      = "goimports -w ."
  dependencies = concat(glob("**.go"), "go.mod", "go.sum")
}
rule {
  target       = ".test"
  tee_target   = true
  dependencies = concat(".import", glob("**/testdata/**"))
  command      = "go test -count=1 ./..."
}
rule {
  target       = ".lint"
  tee_target   = true
  dependencies = ".test"
  command      = "golangci-lint run --fix"
}
