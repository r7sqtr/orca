package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/app"
	"github.com/r7sqtr/orca/internal/config"
	"github.com/r7sqtr/orca/internal/docker"
	"github.com/r7sqtr/orca/internal/i18n"
)

// アプリケーションを実行
func Execute() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "設定ファイルの読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// 言語設定
	i18n.SetLanguage(cfg.Language)

	// dockerバイナリパスを解決
	dockerPath, err := docker.ResolveDockerPath(cfg.DockerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// アプリケーションモデルの作成
	model := app.NewAppModel(cfg, dockerPath)

	// bubbletea プログラムの作成・実行
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
