package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// 既知のdockerバイナリパス候補 (macOS/Linux)
var defaultDockerPaths = []string{
	"/usr/local/bin/docker",
	"/opt/homebrew/bin/docker",
	"/usr/bin/docker",
	"/snap/bin/docker",
}

// ResolveDockerPath はdockerバイナリの絶対パスを解決する。
// 優先順位: configPath(設定ファイル) > exec.LookPath > 既知パスのフォールバック
func ResolveDockerPath(configPath string) (string, error) {
	// 設定ファイルで明示指定されている場合
	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			return "", fmt.Errorf("設定の docker_path が見つかりません: %s: %w", configPath, err)
		}
		return configPath, nil
	}

	// PATHから検索
	if p, err := exec.LookPath("docker"); err == nil {
		return p, nil
	}

	// 既知パスのフォールバック
	for _, p := range defaultDockerPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf(
		"docker が見つかりません。PATH に docker を追加するか、設定ファイル (~/.config/orca/config.yml) で docker_path を指定してください",
	)
}
