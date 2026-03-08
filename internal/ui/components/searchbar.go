package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/ui"
)

// SearchBar は検索入力バー
type SearchBar struct {
	styles ui.Styles
	input  textinput.Model
	active bool
	width  int
}

// NewSearchBar はSearchBarを作成する
func NewSearchBar(styles ui.Styles) SearchBar {
	ti := textinput.New()
	ti.Placeholder = i18n.T("search.placeholder")
	ti.CharLimit = 100

	return SearchBar{
		styles: styles,
		input:  ti,
	}
}

// SetSize はサイズを設定する
func (sb *SearchBar) SetSize(width int) {
	sb.width = width
	sb.input.Width = width - 4
}

// Activate は検索バーを有効化する
func (sb *SearchBar) Activate() {
	sb.active = true
	sb.input.Focus()
}

// Deactivate は検索バーを無効化する
func (sb *SearchBar) Deactivate() {
	sb.active = false
	sb.input.Blur()
}

// Reset は検索をリセットする
func (sb *SearchBar) Reset() {
	sb.input.SetValue("")
	sb.Deactivate()
}

// IsActive は検索バーが有効かを返す
func (sb SearchBar) IsActive() bool {
	return sb.active
}

// Query は現在の検索クエリを返す
func (sb SearchBar) Query() string {
	return sb.input.Value()
}

// Update はキー入力を処理する
func (sb *SearchBar) Update(msg tea.Msg) tea.Cmd {
	if !sb.active {
		return nil
	}

	var cmd tea.Cmd
	sb.input, cmd = sb.input.Update(msg)
	return cmd
}

// View は検索バーを描画する
func (sb SearchBar) View() string {
	if !sb.active {
		return ""
	}
	prefix := sb.styles.Subtitle.Render("/")
	return prefix + sb.input.View()
}
