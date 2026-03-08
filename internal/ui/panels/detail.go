package panels

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/model"
	"github.com/vvsaito/orca/internal/ui"
)

// DetailTab はタブ種別
type DetailTab int

const (
	TabInfo DetailTab = iota
	TabLogs
	TabEnv
)

// Detail は右側の詳細パネル
type Detail struct {
	styles      ui.Styles
	keymap      ui.KeyMap
	activeTab   DetailTab
	serviceInfo ServiceInfo
	logView     LogView
	envView     EnvView
	focused     bool
	width       int
	height      int
}

// NewDetail はDetailパネルを作成する
func NewDetail(styles ui.Styles, keymap ui.KeyMap, logBufferSize int) Detail {
	return Detail{
		styles:      styles,
		keymap:      keymap,
		serviceInfo: NewServiceInfo(styles),
		logView:     NewLogView(styles, keymap, logBufferSize),
		envView:     NewEnvView(styles),
		activeTab:   TabInfo,
	}
}

// SetSize はサイズを設定する
func (d *Detail) SetSize(width, height int) {
	d.width = width
	d.height = height

	tabHeight := 1
	contentHeight := height - tabHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	d.serviceInfo.SetSize(width, contentHeight)
	d.logView.SetSize(width, contentHeight)
	d.envView.SetSize(width, contentHeight)
}

// SetFocused はフォーカス状態を設定する
func (d *Detail) SetFocused(focused bool) {
	d.focused = focused
	d.updateChildFocus()
}

// SetService は表示するサービスを設定する
func (d *Detail) SetService(project, service string, container *model.ContainerStatus) {
	d.serviceInfo.SetService(project, service, container)
}

// LogView はLogViewへの参照を返す
func (d *Detail) LogView() *LogView {
	return &d.logView
}

// EnvView はEnvViewへの参照を返す
func (d *Detail) EnvView() *EnvView {
	return &d.envView
}

// ActiveTab はアクティブなタブを返す
func (d Detail) ActiveTab() DetailTab {
	return d.activeTab
}

// SwitchTab はタブを切り替える
func (d *Detail) SwitchTab(tab DetailTab) {
	d.activeTab = tab
	d.updateChildFocus()
}

func (d *Detail) updateChildFocus() {
	d.logView.SetFocused(d.focused && d.activeTab == TabLogs)
	d.envView.SetFocused(d.focused && d.activeTab == TabEnv)
}

// Update はキー入力を処理する
func (d *Detail) Update(msg tea.Msg) tea.Cmd {
	if !d.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// タブ切替は常に受け付ける（検索中は除外＝MainScreenで処理済み）
		switch {
		case key.Matches(msg, d.keymap.Info):
			d.SwitchTab(TabInfo)
			return nil
		case key.Matches(msg, d.keymap.Logs):
			d.SwitchTab(TabLogs)
			return nil
		case key.Matches(msg, d.keymap.EnvVars):
			d.SwitchTab(TabEnv)
			return nil
		}
	}

	// アクティブなタブに委譲
	switch d.activeTab {
	case TabLogs:
		return d.logView.Update(msg)
	case TabEnv:
		return d.envView.Update(msg)
	}

	return nil
}

// View は詳細パネルを描画する
func (d Detail) View() string {
	tabs := d.renderTabs()

	var content string
	switch d.activeTab {
	case TabInfo:
		content = d.serviceInfo.View()
	case TabLogs:
		content = d.logView.View()
	case TabEnv:
		content = d.envView.View()
	}

	return tabs + "\n" + content
}

func (d Detail) renderTabs() string {
	infoTab := i18n.T("detail.tab.info")
	logsTab := i18n.T("detail.tab.logs")
	envTab := i18n.T("detail.tab.env")

	tabNames := []struct {
		text string
		tab  DetailTab
	}{
		{infoTab, TabInfo},
		{logsTab, TabLogs},
		{envTab, TabEnv},
	}

	var rendered []string
	for _, t := range tabNames {
		if d.activeTab == t.tab {
			rendered = append(rendered, d.styles.ActiveTab.Render(t.text))
		} else {
			rendered = append(rendered, d.styles.InactiveTab.Render(t.text))
		}
	}

	// フォーカスインジケータ
	focusIndicator := ""
	if d.focused {
		focusIndicator = d.styles.Running.Render(" ◀")
	}

	result := ""
	for i, r := range rendered {
		if i > 0 {
			result += " "
		}
		result += r
	}

	return result + focusIndicator
}
