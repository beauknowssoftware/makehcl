package definition

type Rule struct {
	IsPhony      bool
	Target       Target
	Command      string
	Dependencies []Target
	Environment  map[string]string
}
