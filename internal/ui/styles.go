package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/r7sqtr/orca/internal/config"
)

// アプリケーション全体のスタイル定義
type Styles struct {
	Theme config.ThemeColors

	// レイアウト
	App       lipgloss.Style
	Header    lipgloss.Style
	Sidebar   lipgloss.Style
	Detail    lipgloss.Style
	StatusBar lipgloss.Style
	HelpBar   lipgloss.Style

	// テキスト
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Muted       lipgloss.Style
	Bold        lipgloss.Style
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style

	// 状態表示
	Running lipgloss.Style
	Stopped lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Health  lipgloss.Style

	// ログ
	LogStdout    lipgloss.Style
	LogStderr    lipgloss.Style
	LogTimestamp lipgloss.Style
	LogService   lipgloss.Style
	LogHighlight lipgloss.Style

	// サイドバーアイテム
	SelectedItem lipgloss.Style
	NormalItem   lipgloss.Style
	ProjectItem  lipgloss.Style

	// ダイアログ
	Dialog      lipgloss.Style
	DialogTitle lipgloss.Style

	// 選択インジケータ
	ActiveBorder   lipgloss.Style
	InactiveBorder lipgloss.Style
}

// テーマに基づくスタイルを作成
func NewStyles(theme config.ThemeColors) Styles {
	s := Styles{Theme: theme}

	s.App = lipgloss.NewStyle()

	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		Padding(0, 1)

	s.Sidebar = lipgloss.NewStyle().
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Border)

	s.Detail = lipgloss.NewStyle().
		Padding(0, 1)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Background(theme.Border).
		Padding(0, 1)

	s.HelpBar = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 1)

	// テキスト
	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(theme.Secondary)

	s.Muted = lipgloss.NewStyle().
		Foreground(theme.Muted)

	s.Bold = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Foreground)

	// タブ（ボーダーなし: 高さを1行に固定）
	s.ActiveTab = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		Underline(true).
		Padding(0, 1)

	s.InactiveTab = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 1)

	// 状態
	s.Running = lipgloss.NewStyle().
		Foreground(theme.Success)

	s.Stopped = lipgloss.NewStyle().
		Foreground(theme.Muted)

	s.Error = lipgloss.NewStyle().
		Foreground(theme.Error)

	s.Warning = lipgloss.NewStyle().
		Foreground(theme.Warning)

	s.Health = lipgloss.NewStyle().
		Foreground(theme.Success)

	// ログ
	s.LogStdout = lipgloss.NewStyle().
		Foreground(theme.Foreground)

	s.LogStderr = lipgloss.NewStyle().
		Foreground(theme.Error)

	s.LogTimestamp = lipgloss.NewStyle().
		Foreground(theme.Muted)

	s.LogService = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Bold(true)

	s.LogHighlight = lipgloss.NewStyle().
		Background(theme.Highlight).
		Foreground(theme.Foreground)

	// サイドバーアイテム
	s.SelectedItem = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Background(theme.Highlight).
		Padding(0, 1)

	s.NormalItem = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Padding(0, 1)

	s.ProjectItem = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		Padding(0, 1)

	// ダイアログ
	s.Dialog = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	s.DialogTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary)

	// ボーダー
	s.ActiveBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Primary)

	s.InactiveBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Border)

	return s
}
