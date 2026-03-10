package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// docker compose CLIを実行
type ComposeExec struct{}

// ComposeExecを作成
func NewComposeExec() *ComposeExec {
	return &ComposeExec{}
}

// docker composeコマンドを実行する共通ヘルパー
func (ce *ComposeExec) runCommand(ctx context.Context, workingDir string, args []string, combineOutput bool) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	if combineOutput {
		return cmd.CombinedOutput()
	}
	return cmd.Output()
}

// docker composeコマンドを実行
func (ce *ComposeExec) Run(ctx context.Context, workingDir string, action ComposeAction, service string) error {
	args := []string{"compose"}

	switch action {
	case ActionUp:
		args = append(args, "up", "-d")
		if service != "" {
			args = append(args, service)
		}
	case ActionStop:
		args = append(args, "stop")
		if service != "" {
			args = append(args, service)
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

	output, err := ce.runCommand(ctx, workingDir, args, true)
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

// プロジェクト全体に対する操作を実行
func (ce *ComposeExec) ProjectAction(ctx context.Context, workingDir string, action ComposeAction) error {
	return ce.Run(ctx, workingDir, action, "")
}

// 特定サービスに対する操作を実行
func (ce *ComposeExec) ServiceAction(ctx context.Context, workingDir string, action ComposeAction, service string) error {
	return ce.Run(ctx, workingDir, action, service)
}

// compose設定ファイルからサービス名一覧を取得
func (ce *ComposeExec) ListServices(ctx context.Context, workingDir string) ([]string, error) {
	output, err := ce.runCommand(ctx, workingDir, []string{"compose", "config", "--services"}, false)
	if err != nil {
		return nil, fmt.Errorf("docker compose config --services: %w", err)
	}

	var services []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			services = append(services, line)
		}
	}
	return services, nil
}

// docker compose exec用のexec.Cmdを返す
func (ce *ComposeExec) ExecCommand(workingDir, service string) *exec.Cmd {
	cmd := exec.Command("docker", "compose", "exec", service, "sh")
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	return cmd
}
