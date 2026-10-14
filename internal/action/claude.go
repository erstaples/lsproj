package action

import (
	"os"
	"os/exec"
)

type ClaudeAction struct{}

func (a *ClaudeAction) Name() string { return "Run `claude`" }

func (a *ClaudeAction) Run(projectPath string) error {
	cmd := exec.Command("claude")
	cmd.Dir = projectPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
