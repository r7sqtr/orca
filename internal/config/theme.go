package config

import "github.com/charmbracelet/lipgloss"

// ThemeColors はテーマの色定義
type ThemeColors struct {
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Muted      lipgloss.Color
	Border     lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Highlight  lipgloss.Color
}

// DarkTheme はダークテーマの色定義
func DarkTheme() ThemeColors {
	return ThemeColors{
		Primary:    lipgloss.Color("#7dcfff"),
		Secondary:  lipgloss.Color("#bb9af7"),
		Accent:     lipgloss.Color("#e0af68"),
		Success:    lipgloss.Color("#9ece6a"),
		Warning:    lipgloss.Color("#ff9e64"),
		Error:      lipgloss.Color("#f7768e"),
		Muted:      lipgloss.Color("#565f89"),
		Border:     lipgloss.Color("#3b4261"),
		Background: lipgloss.Color("#1a1b26"),
		Foreground: lipgloss.Color("#c0caf5"),
		Highlight:  lipgloss.Color("#33467c"),
	}
}

// LightTheme はライトテーマの色定義
func LightTheme() ThemeColors {
	return ThemeColors{
		Primary:    lipgloss.Color("#2e7de9"),
		Secondary:  lipgloss.Color("#7847bd"),
		Accent:     lipgloss.Color("#8c6c3e"),
		Success:    lipgloss.Color("#587539"),
		Warning:    lipgloss.Color("#b15c00"),
		Error:      lipgloss.Color("#c64343"),
		Muted:      lipgloss.Color("#8990b3"),
		Border:     lipgloss.Color("#c4c8da"),
		Background: lipgloss.Color("#d5d6db"),
		Foreground: lipgloss.Color("#3760bf"),
		Highlight:  lipgloss.Color("#b6bfe2"),
	}
}

// GetTheme はテーマ名に対応する色定義を返す
func GetTheme(name string) ThemeColors {
	switch name {
	case "light":
		return LightTheme()
	default:
		return DarkTheme()
	}
}
