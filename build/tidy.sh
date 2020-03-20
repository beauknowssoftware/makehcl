go mod tidy
find . -name *.hcl | xargs ./tools/hclfmt -w