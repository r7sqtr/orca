package model

// AppConfig はアプリケーション設定
type AppConfig struct {
	Language       string         `yaml:"language"`
	Theme          string         `yaml:"theme"` // "dark", "light", "auto"
	LogBufferSize  int            `yaml:"log_buffer_size"`
	DockerHost     string         `yaml:"docker_host"`
	KeyBindings    map[string]string `yaml:"keybindings"`
	SidebarWidth   int            `yaml:"sidebar_width"` // パーセント (0で自動)
	ConfirmActions bool           `yaml:"confirm_actions"`
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig() AppConfig {
	return AppConfig{
		Language:       "ja",
		Theme:          "dark",
		LogBufferSize:  10000,
		SidebarWidth:   0,
		ConfirmActions: true,
	}
}
