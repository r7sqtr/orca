package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/ui"
)

// 画面下部のステータスバー
type StatusBar struct {
	styles    ui.Styles
	width     int
	connected bool
	message   string
	project   string
	service   string
}

// StatusBarを作成
func NewStatusBar(styles ui.Styles) StatusBar {
	return StatusBar{styles: styles}
}

// サイズを設定
func (sb *StatusBar) SetSize(width int) {
	sb.width = width
}

// 接続状態を設定
func (sb *StatusBar) SetConnected(connected bool) {
	sb.connected = connected
}

// メッセージを設定
func (sb *StatusBar) SetMessage(msg string) {
	sb.message = msg
}

// 現在のコンテキストを設定
func (sb *StatusBar) SetContext(project, service string) {
	sb.project = project
	sb.service = service
}

// ステータスバーを描画
func (sb StatusBar) View() string {
	left := ""
	if sb.connected {
		left = sb.styles.Running.Render("●") + " " + i18n.T("app.connected")
	} else {
		left = sb.styles.Error.Render("●") + " " + i18n.T("app.no_docker")
	}

	if sb.project != "" {
		ctx := sb.project
		if sb.service != "" {
			ctx = fmt.Sprintf("%s/%s", sb.project, sb.service)
		}
		left += " │ " + ctx
	}

	right := sb.message

	// 右揃え
	gap := sb.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	bar := left + lipgloss.NewStyle().Width(gap).Render("") + right
	return sb.styles.StatusBar.Width(sb.width).Render(bar)
}
