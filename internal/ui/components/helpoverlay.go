package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/ui"
)

// キーバインド一覧のオーバーレイ
type HelpOverlay struct {
	styles ui.Styles
	active bool
}

// HelpOverlayを作成
func NewHelpOverlay(styles ui.Styles) HelpOverlay {
	return HelpOverlay{styles: styles}
}

// オーバーレイを表示
func (ho *HelpOverlay) Show() {
	ho.active = true
}

// オーバーレイを非表示に
func (ho *HelpOverlay) Hide() {
	ho.active = false
}

// オーバーレイが表示中かを返す
func (ho HelpOverlay) IsActive() bool {
	return ho.active
}

// キー入力を処理
func (ho *HelpOverlay) Update(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "esc", "?", "q":
		ho.active = false
		return true
	}
	return false
}

// キーバインドの1行分
type keyEntry struct {
	key  string
	desc string
}

// オーバーレイを描画
func (ho HelpOverlay) View(width, height int) string {
	if !ho.active {
		return ""
	}

	title := ho.styles.DialogTitle.Render(i18n.T("help.overlay.title"))

	sidebarEntries := []keyEntry{
		{"j/k", i18n.T("help.desc.move")},
		{"u", i18n.T("help.desc.up")},
		{"d", i18n.T("help.desc.down")},
		{"r", i18n.T("help.desc.restart")},
		{"b", i18n.T("help.desc.build")},
		{"e", i18n.T("help.desc.exec")},
		{"i/l/v", i18n.T("help.desc.tab_switch")},
		{"h", i18n.T("help.desc.toggle")},
		{"Tab", i18n.T("help.desc.panel_switch")},
	}

	detailEntries := []keyEntry{
		{"j/k", i18n.T("help.desc.move")},
		{"i/l/v", i18n.T("help.desc.tab_switch")},
		{"f", i18n.T("help.desc.follow")},
		{"/", i18n.T("help.desc.search")},
		{"y", i18n.T("help.desc.copy")},
		{"o", i18n.T("help.desc.export")},
		{"Esc", i18n.T("help.desc.back")},
	}

	globalEntries := []keyEntry{
		{"?", i18n.T("help.desc.help")},
		{"q", i18n.T("help.desc.quit")},
	}

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	// サイドバーセクション
	lines = append(lines, ho.styles.Subtitle.Render(i18n.T("help.overlay.sidebar")))
	lines = append(lines, ho.formatEntries(sidebarEntries)...)
	lines = append(lines, "")

	// Detailセクション
	lines = append(lines, ho.styles.Subtitle.Render(i18n.T("help.overlay.detail")))
	lines = append(lines, ho.formatEntries(detailEntries)...)
	lines = append(lines, "")

	// 共通セクション
	lines = append(lines, ho.styles.Subtitle.Render(i18n.T("help.overlay.global")))
	lines = append(lines, ho.formatEntries(globalEntries)...)
	lines = append(lines, "")

	lines = append(lines, ho.styles.Muted.Render(i18n.T("help.overlay.close")))

	content := strings.Join(lines, "\n")
	dialog := ho.styles.Dialog.Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, dialog)
}

func (ho HelpOverlay) formatEntries(entries []keyEntry) []string {
	var lines []string
	for _, e := range entries {
		keyStyled := ho.styles.Bold.Render(fmt.Sprintf("  %-8s", e.key))
		lines = append(lines, keyStyled+" "+e.desc)
	}
	return lines
}
