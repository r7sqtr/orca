package panels

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/ui"
)

// 環境変数表示パネル
type EnvView struct {
	styles   ui.Styles
	viewport viewport.Model
	envVars  []string
	focused  bool
	width    int
	height   int
}

// EnvViewを作成
func NewEnvView(styles ui.Styles) EnvView {
	vp := viewport.New(80, 20)
	return EnvView{
		styles:   styles,
		viewport: vp,
	}
}

// サイズを設定
func (ev *EnvView) SetSize(width, height int) {
	ev.width = width
	ev.height = height

	headerHeight := 1
	vpHeight := height - headerHeight
	if vpHeight < 0 {
		vpHeight = 0
	}
	ev.viewport.Width = width
	ev.viewport.Height = vpHeight
}

// フォーカス状態を設定
func (ev *EnvView) SetFocused(focused bool) {
	ev.focused = focused
}

// 環境変数を設定
func (ev *EnvView) SetEnvVars(vars []string) {
	ev.envVars = vars
	ev.refreshContent()
}

// 環境変数をクリア
func (ev *EnvView) Clear() {
	ev.envVars = nil
	ev.refreshContent()
}

// キー入力を処理
func (ev *EnvView) Update(msg tea.Msg) tea.Cmd {
	if !ev.focused {
		return nil
	}
	var cmd tea.Cmd
	ev.viewport, cmd = ev.viewport.Update(msg)
	return cmd
}

// 環境変数パネルを描画
func (ev EnvView) View() string {
	header := ev.styles.Subtitle.Render(i18n.T("env.title"))
	return header + "\n" + ev.viewport.View()
}

func (ev *EnvView) refreshContent() {
	if len(ev.envVars) == 0 {
		ev.viewport.SetContent(ev.styles.Muted.Render(i18n.T("env.no_env")))
		return
	}

	// ソートして表示
	sorted := make([]string, len(ev.envVars))
	copy(sorted, ev.envVars)
	sort.Strings(sorted)

	var lines []string
	for _, env := range sorted {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			key := ev.styles.Bold.Render(parts[0])
			val := ev.styles.Muted.Render("=" + parts[1])
			lines = append(lines, key+val)
		} else {
			lines = append(lines, env)
		}
	}

	ev.viewport.SetContent(strings.Join(lines, "\n"))
}
