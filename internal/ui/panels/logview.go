package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/model"
	"github.com/r7sqtr/orca/internal/ui"
	"github.com/r7sqtr/orca/internal/ui/components"
)

// ログビューアパネル
type LogView struct {
	styles    ui.Styles
	keymap    ui.KeyMap
	viewport  viewport.Model
	searchbar components.SearchBar
	buffer    *model.RingBuffer
	filter    model.LogFilter
	follow    bool
	focused   bool
	width     int
	height    int
}

// LogViewを作成
func NewLogView(styles ui.Styles, keymap ui.KeyMap, bufferSize int) LogView {
	vp := viewport.New(80, 20)

	return LogView{
		styles:    styles,
		keymap:    keymap,
		viewport:  vp,
		searchbar: components.NewSearchBar(styles),
		buffer:    model.NewRingBuffer(bufferSize),
		filter:    model.DefaultLogFilter(),
		follow:    true,
	}
}

// サイズを設定
func (lv *LogView) SetSize(width, height int) {
	lv.width = width
	lv.height = height

	searchHeight := 0
	if lv.searchbar.IsActive() {
		searchHeight = 1
	}
	headerHeight := 1

	vpHeight := height - headerHeight - searchHeight
	if vpHeight < 0 {
		vpHeight = 0
	}

	lv.viewport.Width = width
	lv.viewport.Height = vpHeight
	lv.searchbar.SetSize(width)
}

// フォーカス状態を設定
func (lv *LogView) SetFocused(focused bool) {
	lv.focused = focused
}

// ログエントリを追加
func (lv *LogView) AddEntry(entry model.LogEntry) {
	lv.buffer.Add(entry)
	lv.refreshContent()

	if lv.follow {
		lv.viewport.GotoBottom()
	}
}

// ログをクリア
func (lv *LogView) Clear() {
	lv.buffer.Clear()
	lv.refreshContent()
}

// フォローモードを切り替える
func (lv *LogView) ToggleFollow() {
	lv.follow = !lv.follow
	if lv.follow {
		lv.viewport.GotoBottom()
	}
}

// フォローモードかを返す
func (lv LogView) IsFollowing() bool {
	return lv.follow
}

// 検索バーがアクティブかを返す
func (lv LogView) IsSearchActive() bool {
	return lv.searchbar.IsActive()
}

// キー入力を処理
func (lv *LogView) Update(msg tea.Msg) tea.Cmd {
	if !lv.focused {
		return nil
	}

	// 検索バーがアクティブなら検索バーに委譲
	if lv.searchbar.IsActive() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				lv.searchbar.Reset()
				lv.filter.SearchQuery = ""
				lv.refreshContent()
				return nil
			case "enter":
				lv.filter.SearchQuery = lv.searchbar.Query()
				lv.searchbar.Deactivate()
				lv.refreshContent()
				return nil
			}
		}
		cmd := lv.searchbar.Update(msg)
		// リアルタイム検索
		lv.filter.SearchQuery = lv.searchbar.Query()
		lv.refreshContent()
		return cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, lv.keymap.Search):
			lv.searchbar.Activate()
			lv.SetSize(lv.width, lv.height) // 高さ再計算
			return nil
		case key.Matches(msg, lv.keymap.Follow):
			lv.ToggleFollow()
			return nil
		}
	}

	// スクロール操作を検出してfollowを無効化
	prevOffset := lv.viewport.YOffset

	var cmd tea.Cmd
	lv.viewport, cmd = lv.viewport.Update(msg)

	// ユーザーが手動スクロールしたらfollowを無効化
	if lv.viewport.YOffset != prevOffset && lv.follow {
		atBottom := lv.viewport.AtBottom()
		if !atBottom {
			lv.follow = false
		}
	}

	return cmd
}

// ログビューを描画
func (lv LogView) View() string {
	// ヘッダー
	header := lv.renderHeader()

	// 検索バー
	searchView := ""
	if lv.searchbar.IsActive() {
		searchView = lv.searchbar.View() + "\n"
	}

	return header + "\n" + searchView + lv.viewport.View()
}

func (lv LogView) renderHeader() string {
	title := lv.styles.Subtitle.Render(i18n.T("log.title"))

	// フォロー状態
	followText := ""
	if lv.follow {
		followText = lv.styles.Running.Render(" " + i18n.T("log.follow"))
	} else {
		followText = lv.styles.Muted.Render(" " + i18n.T("log.paused"))
	}

	// 行数
	countText := lv.styles.Muted.Render(
		fmt.Sprintf(" %s", i18n.TF("log.lines", lv.buffer.Count())),
	)

	// 検索結果
	searchText := ""
	if lv.filter.SearchQuery != "" {
		entries := lv.buffer.FilteredEntries(lv.filter)
		searchText = lv.styles.Warning.Render(
			fmt.Sprintf(" %s", i18n.TF("log.matches", len(entries))),
		)
	}

	return title + followText + countText + searchText
}

func (lv *LogView) refreshContent() {
	entries := lv.buffer.FilteredEntries(lv.filter)

	if len(entries) == 0 {
		lv.viewport.SetContent(lv.styles.Muted.Render(i18n.T("log.no_logs")))
		return
	}

	var lines []string
	for _, entry := range entries {
		line := lv.formatEntry(entry)
		lines = append(lines, line)
	}

	lv.viewport.SetContent(strings.Join(lines, "\n"))
}

func (lv LogView) formatEntry(entry model.LogEntry) string {
	// タイムスタンプ
	ts := lv.styles.LogTimestamp.Render(entry.Timestamp.Format("15:04:05"))

	// サービス名
	svc := lv.styles.LogService.Render(entry.Service)

	// メッセージ
	msg := entry.Message
	if lv.filter.SearchQuery != "" {
		msg = lv.highlightSearch(msg, lv.filter.SearchQuery)
	}

	var msgStyled string
	if entry.Stream == model.StreamStderr {
		msgStyled = lv.styles.LogStderr.Render(msg)
	} else {
		msgStyled = lv.styles.LogStdout.Render(msg)
	}

	return fmt.Sprintf("%s %s %s", ts, svc, msgStyled)
}

// ログのプレーンテキストを返す（コピー・エクスポート用）
func (lv LogView) GetPlainText() string {
	entries := lv.buffer.FilteredEntries(lv.filter)
	if len(entries) == 0 {
		return ""
	}

	var lines []string
	for _, entry := range entries {
		line := fmt.Sprintf("%s %s %s",
			entry.Timestamp.Format("15:04:05"),
			entry.Service,
			entry.Message,
		)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (lv LogView) highlightSearch(text, query string) string {
	lower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	var result strings.Builder
	pos := 0
	for {
		idx := strings.Index(lower[pos:], queryLower)
		if idx < 0 {
			result.WriteString(text[pos:])
			break
		}
		result.WriteString(text[pos : pos+idx])
		matchEnd := pos + idx + len(query)
		result.WriteString(lv.styles.LogHighlight.Render(text[pos+idx : matchEnd]))
		pos = matchEnd
	}
	return result.String()
}
