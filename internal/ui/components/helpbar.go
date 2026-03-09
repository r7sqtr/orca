package components

import (
	"strings"

	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/ui"
)

// HelpBar はキーヒントを表示するバー
type HelpBar struct {
	styles ui.Styles
	width  int
	mode   HelpMode
}

// HelpMode はヘルプ表示モード
type HelpMode int

const (
	HelpModeNormal  HelpMode = iota // サイドバーフォーカス
	HelpModeSearch                  // 検索入力中
	HelpModeConfirm                 // 確認ダイアログ
	HelpModeLogs                    // Detailフォーカス + ログタブ
	HelpModeInfo                    // Detailフォーカス + 情報タブ
	HelpModeEnv                     // Detailフォーカス + 環境変数タブ
)

// NewHelpBar はHelpBarを作成する
func NewHelpBar(styles ui.Styles) HelpBar {
	return HelpBar{styles: styles}
}

// SetSize はサイズを設定する
func (hb *HelpBar) SetSize(width int) {
	hb.width = width
}

// SetMode はヘルプモードを設定する
func (hb *HelpBar) SetMode(mode HelpMode) {
	hb.mode = mode
}

// View はヘルプバーを描画する
func (hb HelpBar) View() string {
	var keys []string

	switch hb.mode {
	case HelpModeSearch:
		keys = []string{
			i18n.T("help.esc"),
			i18n.T("help.enter"),
		}
	case HelpModeConfirm:
		keys = []string{
			"[y]" + i18n.T("confirm.yes"),
			"[n]" + i18n.T("confirm.no"),
		}
	case HelpModeLogs:
		keys = []string{
			i18n.T("help.move"),
			i18n.T("help.follow"),
			i18n.T("help.search"),
			i18n.T("help.copy"),
			i18n.T("help.export"),
			i18n.T("help.tab"),
			i18n.T("help.focus"),
			i18n.T("help.esc"),
		}
	case HelpModeInfo:
		keys = []string{
			i18n.T("help.tab"),
			i18n.T("help.focus"),
			i18n.T("help.esc"),
		}
	case HelpModeEnv:
		keys = []string{
			i18n.T("help.move"),
			i18n.T("help.tab"),
			i18n.T("help.focus"),
			i18n.T("help.esc"),
		}
	default:
		// サイドバーフォーカス
		keys = []string{
			i18n.T("help.move"),
			i18n.T("help.up"),
			i18n.T("help.down"),
			i18n.T("help.restart"),
			i18n.T("help.build"),
			i18n.T("help.exec"),
			i18n.T("help.info"),
			i18n.T("help.logs"),
			i18n.T("help.env"),
			i18n.T("help.focus"),
			i18n.T("help.help"),
			i18n.T("help.quit"),
		}
	}

	text := strings.Join(keys, " ")
	return hb.styles.HelpBar.Width(hb.width).Render(text)
}
