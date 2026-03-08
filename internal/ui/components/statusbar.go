package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/ui"
)

// StatusBar は画面下部のステータスバー
type StatusBar struct {
	styles    ui.Styles
	width     int
	connected bool
	message   string
	project   string
	service   string
}

// NewStatusBar はStatusBarを作成する
func NewStatusBar(styles ui.Styles) StatusBar {
	return StatusBar{styles: styles}
}

// SetSize はサイズを設定する
func (sb *StatusBar) SetSize(width int) {
	sb.width = width
}

// SetConnected は接続状態を設定する
func (sb *StatusBar) SetConnected(connected bool) {
	sb.connected = connected
}

// SetMessage はメッセージを設定する
func (sb *StatusBar) SetMessage(msg string) {
	sb.message = msg
}

// SetContext は現在のコンテキストを設定する
func (sb *StatusBar) SetContext(project, service string) {
	sb.project = project
	sb.service = service
}

// View はステータスバーを描画する
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
