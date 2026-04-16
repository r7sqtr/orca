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
)

// ボリューム一覧表示パネル
type VolumeView struct {
	styles   ui.Styles
	keymap   ui.KeyMap
	viewport viewport.Model
	volumes  []model.VolumeInfo
	cursor   int
	focused  bool
	width    int
	height   int
}

// VolumeViewを作成
func NewVolumeView(styles ui.Styles, keymap ui.KeyMap) VolumeView {
	vp := viewport.New(80, 20)
	return VolumeView{
		styles:   styles,
		keymap:   keymap,
		viewport: vp,
	}
}

// サイズを設定
func (vv *VolumeView) SetSize(width, height int) {
	vv.width = width
	vv.height = height

	headerHeight := 2 // タイトル + テーブルヘッダー
	vpHeight := height - headerHeight
	if vpHeight < 0 {
		vpHeight = 0
	}
	vv.viewport.Width = width
	vv.viewport.Height = vpHeight
}

// フォーカス状態を設定
func (vv *VolumeView) SetFocused(focused bool) {
	vv.focused = focused
}

// ボリューム一覧を設定
func (vv *VolumeView) SetVolumes(volumes []model.VolumeInfo) {
	vv.volumes = volumes
	if vv.cursor >= len(volumes) {
		vv.cursor = len(volumes) - 1
	}
	if vv.cursor < 0 {
		vv.cursor = 0
	}
	vv.refreshContent()
}

// クリア
func (vv *VolumeView) Clear() {
	vv.volumes = nil
	vv.cursor = 0
	vv.refreshContent()
}

// 選択中のボリュームを返す
func (vv *VolumeView) SelectedVolume() *model.VolumeInfo {
	if vv.cursor >= 0 && vv.cursor < len(vv.volumes) {
		vol := vv.volumes[vv.cursor]
		return &vol
	}
	return nil
}

// キー入力を処理
func (vv *VolumeView) Update(msg tea.Msg) tea.Cmd {
	if !vv.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, vv.keymap.Up):
			if vv.cursor > 0 {
				vv.cursor--
				vv.refreshContent()
				vv.ensureVisible()
			}
			return nil
		case key.Matches(msg, vv.keymap.Down):
			if vv.cursor < len(vv.volumes)-1 {
				vv.cursor++
				vv.refreshContent()
				vv.ensureVisible()
			}
			return nil
		}
	}

	var cmd tea.Cmd
	vv.viewport, cmd = vv.viewport.Update(msg)
	return cmd
}

// スクロール位置を調整してカーソルが見えるようにする
func (vv *VolumeView) ensureVisible() {
	vpHeight := vv.viewport.Height
	if vpHeight < 1 {
		return
	}
	if vv.cursor < vv.viewport.YOffset {
		vv.viewport.SetYOffset(vv.cursor)
	}
	if vv.cursor >= vv.viewport.YOffset+vpHeight {
		vv.viewport.SetYOffset(vv.cursor - vpHeight + 1)
	}
}

// ボリューム一覧パネルを描画
func (vv VolumeView) View() string {
	header := vv.styles.Subtitle.Render(i18n.T("volumes.title"))

	if len(vv.volumes) == 0 {
		return header + "\n" + vv.styles.Muted.Render(i18n.T("volumes.no_volumes"))
	}

	// テーブルヘッダー
	tableHeader := vv.styles.Bold.Render(vv.formatRow(
		i18n.T("volumes.name"),
		i18n.T("volumes.driver"),
		i18n.T("volumes.status"),
	))

	return header + "\n" + tableHeader + "\n" + vv.viewport.View()
}

func (vv *VolumeView) refreshContent() {
	if len(vv.volumes) == 0 {
		vv.viewport.SetContent("")
		return
	}

	var lines []string
	for idx, vol := range vv.volumes {
		name := vol.Name
		if len(name) > 30 {
			name = name[:12] + "..." + name[len(name)-12:]
		}
		driver := vol.Driver

		var status string
		if vol.IsUnused() {
			status = vv.styles.Muted.Render(i18n.T("volumes.unused"))
		} else {
			status = vv.styles.Running.Render(i18n.T("volumes.used"))
		}

		row := vv.formatRow(name, driver, status)

		if idx == vv.cursor && vv.focused {
			lines = append(lines, vv.styles.SelectedItem.Width(vv.width-1).Render(row))
		} else {
			lines = append(lines, vv.styles.NormalItem.Width(vv.width-1).Render(row))
		}
	}

	vv.viewport.SetContent(strings.Join(lines, "\n"))
}

func (vv VolumeView) formatRow(col1, col2, col3 string) string {
	w := vv.width
	if w < 40 {
		w = 40
	}
	driverW := 10
	statusW := 10
	nameW := w - driverW - statusW - 4
	if nameW < 10 {
		nameW = 10
	}

	return fmt.Sprintf("%-*s  %-*s  %s", nameW, col1, driverW, col2, col3)
}
