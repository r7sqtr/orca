package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vvsaito/orca/internal/app"
	"github.com/vvsaito/orca/internal/config"
	"github.com/vvsaito/orca/internal/i18n"
)

// Execute はアプリケーションを実行する
func Execute() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "設定ファイルの読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// 言語設定
	i18n.SetLanguage(cfg.Language)

	// アプリケーションモデルの作成
	model := app.NewAppModel(cfg)

	// bubbletea プログラムの作成・実行
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
