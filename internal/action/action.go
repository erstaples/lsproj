package action

type Action interface {
	Name() string
	Run(projectPath string) error
}

var registry []Action

func Register(a Action) {
	registry = append(registry, a)
}

func All() []Action {
	return registry
}

func init() {
	Register(&CDAction{})
	Register(&ClaudeAction{})
	Register(&ClaudeResumeAction{})
	Register(&VSCodeAction{})
}

// IsCDAction returns true if the action is the special cd action
// that needs shell wrapper handling.
func IsCDAction(a Action) bool {
	_, ok := a.(*CDAction)
	return ok
}
