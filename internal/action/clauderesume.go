package action

import (
	"os"
	"os/exec"
)

type ClaudeResumeAction struct{}

func (a *ClaudeResumeAction) Name() string { return "Run `claude --resume`" }

func (a *ClaudeResumeAction) Run(projectPath string) error {
	cmd := exec.Command("claude", "--resume")
	cmd.Dir = projectPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
