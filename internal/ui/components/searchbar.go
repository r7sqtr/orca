package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/ui"
)

// 検索入力バー
type SearchBar struct {
	styles ui.Styles
	input  textinput.Model
	active bool
	width  int
}

// SearchBarを作成
func NewSearchBar(styles ui.Styles) SearchBar {
	ti := textinput.New()
	ti.Placeholder = i18n.T("search.placeholder")
	ti.CharLimit = 100

	return SearchBar{
		styles: styles,
		input:  ti,
	}
}

// サイズを設定
func (sb *SearchBar) SetSize(width int) {
	sb.width = width
	sb.input.Width = width - 4
}

// 検索バーを有効化
func (sb *SearchBar) Activate() {
	sb.active = true
	sb.input.Focus()
}

// 検索バーを無効化
func (sb *SearchBar) Deactivate() {
	sb.active = false
	sb.input.Blur()
}

// 検索をリセット
func (sb *SearchBar) Reset() {
	sb.input.SetValue("")
	sb.Deactivate()
}

// 検索バーが有効かを返す
func (sb SearchBar) IsActive() bool {
	return sb.active
}

// 現在の検索クエリを返す
func (sb SearchBar) Query() string {
	return sb.input.Value()
}

// キー入力を処理
func (sb *SearchBar) Update(msg tea.Msg) tea.Cmd {
	if !sb.active {
		return nil
	}

	var cmd tea.Cmd
	sb.input, cmd = sb.input.Update(msg)
	return cmd
}

// 検索バーを描画
func (sb SearchBar) View() string {
	if !sb.active {
		return ""
	}
	prefix := sb.styles.Subtitle.Render("/")
	return prefix + sb.input.View()
}
