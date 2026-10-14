package action

// CDAction is handled specially — it doesn't execute anything itself.
// The program writes the selected path to a result file for the shell wrapper.
type CDAction struct{}

func (a *CDAction) Name() string { return "Run `cd`" }

func (a *CDAction) Run(projectPath string) error {
	// The actual cd is handled by the shell wrapper.
	return nil
}
