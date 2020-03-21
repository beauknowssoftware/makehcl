package graph

type Type int

const (
	ImportGraph Type = iota
)

type Options struct {
	GraphType Type
}

type DoOptions struct {
	Filename           string
	IgnoreParserErrors bool
	Options
}
