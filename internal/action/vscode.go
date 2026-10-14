package action

import "os/exec"

type VSCodeAction struct{}

func (a *VSCodeAction) Name() string { return "Open in VS Code" }

func (a *VSCodeAction) Run(projectPath string) error {
	cmd := exec.Command("code", projectPath)
	return cmd.Start()
}
