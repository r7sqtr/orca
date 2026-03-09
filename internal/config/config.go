package config

import (
	"os"
	"path/filepath"

	"github.com/vvsaito/orca/internal/model"
	"gopkg.in/yaml.v3"
)

// configDir は設定ディレクトリのパスを返す
func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "orca")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "orca")
}

// ConfigPath は設定ファイルのパスを返す
func ConfigPath() string {
	return filepath.Join(configDir(), "config.yml")
}

// Load は設定ファイルを読み込む。存在しない場合はデフォルト設定を返す
func Load() (model.AppConfig, error) {
	cfg := model.DefaultConfig()

	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			// デフォルト設定ファイルを作成（書き込み失敗は無視）
			_ = writeDefaultConfig()
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return model.DefaultConfig(), err
	}

	// バリデーション
	if cfg.LogBufferSize <= 0 {
		cfg.LogBufferSize = 10000
	}
	if cfg.LogBufferSize > 100000 {
		cfg.LogBufferSize = 100000
	}

	return cfg, nil
}

// defaultConfigTemplate はコメント付きのデフォルト設定テンプレート
const defaultConfigTemplate = `# Orca 設定ファイル

# 言語設定: "ja" (日本語), "en" (English)
language: ja

# テーマ: "dark", "light", "auto"
theme: dark

# ログバッファサイズ (1〜100000)
log_buffer_size: 10000

# Docker ホスト (未設定時は環境変数 DOCKER_HOST または自動検出)
# docker_host: ""

# サイドバー幅 (パーセント, 0で自動)
sidebar_width: 0

# 破壊的操作の確認ダイアログ
confirm_actions:
  exec: true
  up: true
  down: true
  restart: true
  build: true

# キーバインドのカスタマイズ
# 各キーにはデフォルト値が設定されています
# keybindings:
#   up: k          # 上へ移動
#   down: j        # 下へ移動
#   select: enter  # 選択
#   back: esc      # 戻る
#   quit: q        # 終了
#   tab: tab       # パネル切替
#   focus_left: ctrl+h   # 左パネルフォーカス
#   focus_right: ctrl+l  # 右パネルフォーカス
#   start: u       # サービス起動
#   stop: d        # サービス停止
#   restart: r     # サービス再起動
#   search: /      # 検索
#   follow: f      # ログフォロー
#   logs: l        # ログ表示
#   info: i        # サービス情報
#   help: "?"      # ヘルプ
#   exec: e        # シェル (exec)
#   copy: y        # コピー
#   export: o      # エクスポート
#   env_vars: v    # 環境変数
#   build: b       # イメージビルド
`

// writeDefaultConfig はデフォルト設定ファイルをコメント付きで書き出す
func writeDefaultConfig() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), []byte(defaultConfigTemplate), 0644)
}

// Save は設定ファイルを保存する
func Save(cfg model.AppConfig) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0644)
}
