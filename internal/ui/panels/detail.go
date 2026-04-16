package panels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/model"
	"github.com/r7sqtr/orca/internal/ui"
)

// タブ種別
type DetailTab int

const (
	TabInfo DetailTab = iota
	TabLogs
	TabEnv
	TabImages
	TabVolumes
	NumDetailTabs // タブの総数
)

// 右側の詳細パネル
type Detail struct {
	styles      ui.Styles
	keymap      ui.KeyMap
	activeTab   DetailTab
	serviceInfo ServiceInfo
	logView     LogView
	envView     EnvView
	imageView   ImageView
	volumeView  VolumeView
	focused     bool
	width       int
	height      int
}

// Detailパネルを作成
func NewDetail(styles ui.Styles, keymap ui.KeyMap, logBufferSize int) Detail {
	return Detail{
		styles:      styles,
		keymap:      keymap,
		serviceInfo: NewServiceInfo(styles),
		logView:     NewLogView(styles, keymap, logBufferSize),
		envView:     NewEnvView(styles),
		imageView:   NewImageView(styles, keymap),
		volumeView:  NewVolumeView(styles, keymap),
		activeTab:   TabInfo,
	}
}

// サイズを設定
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
	d.imageView.SetSize(width, contentHeight)
	d.volumeView.SetSize(width, contentHeight)
}

// フォーカス状態を設定
func (d *Detail) SetFocused(focused bool) {
	d.focused = focused
	d.updateChildFocus()
}

// 表示するサービスを設定
func (d *Detail) SetService(project, service string, container *model.ContainerStatus) {
	d.serviceInfo.SetService(project, service, container)
}

// LogViewへの参照を返す
func (d *Detail) LogView() *LogView {
	return &d.logView
}

// EnvViewへの参照を返す
func (d *Detail) EnvView() *EnvView {
	return &d.envView
}

// ImageViewへの参照を返す
func (d *Detail) ImageView() *ImageView {
	return &d.imageView
}

// VolumeViewへの参照を返す
func (d *Detail) VolumeView() *VolumeView {
	return &d.volumeView
}

// アクティブなタブを返す
func (d Detail) ActiveTab() DetailTab {
	return d.activeTab
}

// タブを切り替える
func (d *Detail) SwitchTab(tab DetailTab) {
	d.activeTab = tab
	d.updateChildFocus()
}

func (d *Detail) updateChildFocus() {
	d.logView.SetFocused(d.focused && d.activeTab == TabLogs)
	d.envView.SetFocused(d.focused && d.activeTab == TabEnv)
	d.imageView.SetFocused(d.focused && d.activeTab == TabImages)
	d.volumeView.SetFocused(d.focused && d.activeTab == TabVolumes)
}

// キー入力を処理
func (d *Detail) Update(msg tea.Msg) tea.Cmd {
	if !d.focused {
		return nil
	}

	// アクティブなタブに委譲（タブ切替はMainScreenで処理済み）
	switch d.activeTab {
	case TabLogs:
		return d.logView.Update(msg)
	case TabEnv:
		return d.envView.Update(msg)
	case TabImages:
		return d.imageView.Update(msg)
	case TabVolumes:
		return d.volumeView.Update(msg)
	}

	return nil
}

// 詳細パネルを描画
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
	case TabImages:
		content = d.imageView.View()
	case TabVolumes:
		content = d.volumeView.View()
	}

	return tabs + "\n" + content
}

func (d Detail) renderTabs() string {
	infoTab := i18n.T("detail.tab.info")
	logsTab := i18n.T("detail.tab.logs")
	envTab := i18n.T("detail.tab.env")

	// 幅が狭い場合はタブ名を省略
	imagesTab := i18n.T("detail.tab.images")
	volumesTab := i18n.T("detail.tab.volumes")
	if d.width < 60 {
		imagesTab = i18n.T("detail.tab.images.short")
		volumesTab = i18n.T("detail.tab.volumes.short")
	}

	tabNames := []struct {
		text string
		tab  DetailTab
	}{
		{infoTab, TabInfo},
		{logsTab, TabLogs},
		{envTab, TabEnv},
		{imagesTab, TabImages},
		{volumesTab, TabVolumes},
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
