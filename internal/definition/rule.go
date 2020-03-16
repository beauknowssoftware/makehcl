package definition

type Rule struct {
	TeeTarget    bool
	IsPhony      bool
	Target       Target
	Command      string
	Dependencies []Target
	Environment  map[string]string
}
