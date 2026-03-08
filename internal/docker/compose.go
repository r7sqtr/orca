package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ComposeExec はdocker compose CLIを実行する
type ComposeExec struct{}

// NewComposeExec はComposeExecを作成する
func NewComposeExec() *ComposeExec {
	return &ComposeExec{}
}

// Run はdocker composeコマンドを実行する
func (ce *ComposeExec) Run(ctx context.Context, workingDir string, action ComposeAction, service string) error {
	args := []string{"compose"}

	switch action {
	case ActionUp:
		args = append(args, "up", "-d")
		if service != "" {
			args = append(args, service)
		}
	case ActionDown:
		args = append(args, "down")
		// downはサービス指定でstopを使う
		if service != "" {
			args = []string{"compose", "stop", service}
		}
	case ActionRestart:
		args = append(args, "restart")
		if service != "" {
			args = append(args, service)
		}
	case ActionBuild:
		args = append(args, "build")
		if service != "" {
			args = append(args, service)
		}
	default:
		return fmt.Errorf("不明なアクション: %s", action)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

// ProjectAction はプロジェクト全体に対する操作を実行する
func (ce *ComposeExec) ProjectAction(ctx context.Context, workingDir string, action ComposeAction) error {
	return ce.Run(ctx, workingDir, action, "")
}

// ServiceAction は特定サービスに対する操作を実行する
func (ce *ComposeExec) ServiceAction(ctx context.Context, workingDir string, action ComposeAction, service string) error {
	return ce.Run(ctx, workingDir, action, service)
}

// ExecCommand はdocker compose exec用のexec.Cmdを返す
func (ce *ComposeExec) ExecCommand(workingDir, service string) *exec.Cmd {
	cmd := exec.Command("docker", "compose", "exec", service, "sh")
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	return cmd
}
