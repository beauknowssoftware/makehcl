package definition

type Command struct {
	Name         Target
	Command      string
	Dependencies []Target
	Environment  map[string]string
}

func (c Command) AsRule() *Rule {
	return &Rule{
		IsPhony:      true,
		Target:       c.Name,
		Command:      c.Command,
		Dependencies: c.Dependencies,
		Environment:  c.Environment,
	}
}
