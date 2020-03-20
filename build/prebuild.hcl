// go prebuild
rule {
  target       = ".import"
  tee_target   = true
  command      = "tools/goimports -w ."
  dependencies = concat(glob("**.go"), "go.mod", "go.sum", "tools/goimports")
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
  dependencies = [".test", "tools/golangci-lint"]
  command      = "tools/golangci-lint run --fix"
}
