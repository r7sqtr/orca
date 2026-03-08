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
