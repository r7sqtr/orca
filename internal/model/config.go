package model

// アクション毎の確認ダイアログ設定
type ConfirmConfig struct {
	Exec         bool `yaml:"exec"`
	Up           bool `yaml:"up"`
	Stop         bool `yaml:"stop"`
	Down         bool `yaml:"down"`
	Restart      bool `yaml:"restart"`
	Build        bool `yaml:"build"`
	RemoveImage  bool `yaml:"remove_image"`
	RemoveVolume bool `yaml:"remove_volume"`
	PruneImages  bool `yaml:"prune_images"`
	PruneVolumes bool `yaml:"prune_volumes"`
}

// デフォルトの確認設定を返す（全て有効）
func DefaultConfirmConfig() ConfirmConfig {
	return ConfirmConfig{
		Exec:         true,
		Up:           true,
		Stop:         true,
		Down:         true,
		Restart:      true,
		Build:        true,
		RemoveImage:  true,
		RemoveVolume: true,
		PruneImages:  true,
		PruneVolumes: true,
	}
}

// アプリケーション設定
type AppConfig struct {
	Language       string            `yaml:"language"`
	Theme          string            `yaml:"theme"` // "dark", "light", "auto"
	LogBufferSize  int               `yaml:"log_buffer_size"`
	DockerHost     string            `yaml:"docker_host"`
	DockerPath     string            `yaml:"docker_path"`
	KeyBindings    map[string]string `yaml:"keybindings"`
	SidebarWidth   int               `yaml:"sidebar_width"` // パーセント (0で自動)
	ConfirmActions ConfirmConfig     `yaml:"confirm_actions"`
}

// デフォルト設定を返す
func DefaultConfig() AppConfig {
	return AppConfig{
		Language:       "ja",
		Theme:          "dark",
		LogBufferSize:  10000,
		SidebarWidth:   0,
		ConfirmActions: DefaultConfirmConfig(),
	}
}
