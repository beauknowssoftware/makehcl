package definition

type Goal = []string
type Target = string

type Definition struct {
	defaultGoal       Goal
	rules             map[Target]*Rule
	firstExecutorGoal Goal
	GlobalEnvironment map[string]string
	Shell             string
	ShellFlags        *string
}

func (d Definition) Rule(t Target) *Rule {
	if d.rules == nil {
		return nil
	}

	return d.rules[t]
}

func (d *Definition) SetDefaultGoal(goal Goal) {
	d.defaultGoal = goal
}

func (d *Definition) AddRule(r *Rule) {
	if d.rules == nil {
		d.rules = make(map[Target]*Rule)
	}
	d.rules[r.Target] = r
	if len(d.firstExecutorGoal) == 0 {
		d.firstExecutorGoal = Goal{r.Target}
	}
}

func (d *Definition) AddCommand(c *Command) {
	d.AddRule(c.AsRule())
}

func (d Definition) Rules() map[Target]*Rule {
	return d.rules
}

func (d *Definition) EffectiveGoal(explicitGoal Goal) Goal {
	if len(explicitGoal) == 0 {
		if len(d.defaultGoal) == 0 {
			return d.firstExecutorGoal
		}
		return d.defaultGoal
	}
	return explicitGoal
}
