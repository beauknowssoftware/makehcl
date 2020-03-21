package parse2

type StopAfterStage int

const (
	StopAfterImports StopAfterStage = iota + 1
)

type Options struct {
	StopAfterStage StopAfterStage
	Filename       string
}
