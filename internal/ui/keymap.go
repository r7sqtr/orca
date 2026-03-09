package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap はアプリケーションのキーバインド定義
type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding
	Back       key.Binding
	Quit       key.Binding
	Tab        key.Binding
	FocusLeft  key.Binding
	FocusRight key.Binding
	Start      key.Binding
	Stop       key.Binding
	Restart    key.Binding
	Search     key.Binding
	Follow     key.Binding
	Logs       key.Binding
	Info       key.Binding
	Help       key.Binding
	Exec       key.Binding
	Copy       key.Binding
	Export     key.Binding
	EnvVars    key.Binding
	Build      key.Binding
}

// DefaultKeyMap はデフォルトのキーバインドを返す
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "上へ"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "下へ"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "選択"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "戻る"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "終了"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "タブ切替"),
		),
		FocusLeft: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("C-h", "左パネル"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("C-l", "右パネル"),
		),
		Start: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "起動"),
		),
		Stop: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "停止"),
		),
		Restart: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "再起動"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "検索"),
		),
		Follow: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "フォロー"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "ログ"),
		),
		Info: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "情報"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "ヘルプ"),
		),
		Exec: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "シェル"),
		),
		Copy: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "コピー"),
		),
		Export: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "エクスポート"),
		),
		EnvVars: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "環境変数"),
		),
		Build: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "ビルド"),
		),
	}
}
