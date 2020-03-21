package graph

type Type string

const (
	ImportGraph Type = "import"
)

type Options struct {
	GraphType Type
}

type DoOptions struct {
	Filename           string
	IgnoreParserErrors bool
	Options
}
