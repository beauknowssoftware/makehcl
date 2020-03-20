go mod tidy
find . -name *.hcl | xargs ./bin/hclfmt -w