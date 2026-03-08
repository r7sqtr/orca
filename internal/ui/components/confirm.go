package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/ui"
)

// ConfirmResult は確認ダイアログの結果
type ConfirmResult struct {
	Confirmed bool
	Action    string
	Target    string
}

// ConfirmDialog は確認ダイアログ
type ConfirmDialog struct {
	styles  ui.Styles
	active  bool
	message string
	action  string
	target  string
}

// NewConfirmDialog はConfirmDialogを作成する
func NewConfirmDialog(styles ui.Styles) ConfirmDialog {
	return ConfirmDialog{styles: styles}
}

// Show はダイアログを表示する
func (cd *ConfirmDialog) Show(message, action, target string) {
	cd.active = true
	cd.message = message
	cd.action = action
	cd.target = target
}

// Hide はダイアログを非表示にする
func (cd *ConfirmDialog) Hide() {
	cd.active = false
}

// IsActive はダイアログが表示中かを返す
func (cd ConfirmDialog) IsActive() bool {
	return cd.active
}

// Update はキー入力を処理する
func (cd *ConfirmDialog) Update(msg tea.KeyMsg) *ConfirmResult {
	switch msg.String() {
	case "y", "Y", "enter":
		cd.active = false
		return &ConfirmResult{
			Confirmed: true,
			Action:    cd.action,
			Target:    cd.target,
		}
	case "n", "N", "esc":
		cd.active = false
		return &ConfirmResult{
			Confirmed: false,
			Action:    cd.action,
			Target:    cd.target,
		}
	}
	return nil
}

// View はダイアログを描画する
func (cd ConfirmDialog) View(width, height int) string {
	if !cd.active {
		return ""
	}

	title := cd.styles.DialogTitle.Render(i18n.T("confirm.title"))
	message := cd.message
	hint := cd.styles.Muted.Render("[y] " + i18n.T("confirm.yes") + "  [n] " + i18n.T("confirm.no"))

	content := title + "\n\n" + message + "\n\n" + hint

	dialog := cd.styles.Dialog.Render(content)

	// 中央配置
	dialogWidth := lipgloss.Width(dialog)
	dialogHeight := lipgloss.Height(dialog)

	x := (width - dialogWidth) / 2
	y := (height - dialogHeight) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, dialog)
}
