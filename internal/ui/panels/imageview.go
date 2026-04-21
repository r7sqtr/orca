package panels

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/model"
	"github.com/r7sqtr/orca/internal/ui"
)

// イメージ一覧表示パネル
type ImageView struct {
	styles   ui.Styles
	keymap   ui.KeyMap
	viewport viewport.Model
	images   []model.ImageInfo
	cursor   int
	focused  bool
	width    int
	height   int
}

// ImageViewを作成
func NewImageView(styles ui.Styles, keymap ui.KeyMap) ImageView {
	vp := viewport.New(80, 20)
	return ImageView{
		styles:   styles,
		keymap:   keymap,
		viewport: vp,
	}
}

// サイズを設定
func (iv *ImageView) SetSize(width, height int) {
	iv.width = width
	iv.height = height

	headerHeight := 2 // タイトル + テーブルヘッダー
	vpHeight := height - headerHeight
	if vpHeight < 0 {
		vpHeight = 0
	}
	iv.viewport.Width = width
	iv.viewport.Height = vpHeight
}

// フォーカス状態を設定
func (iv *ImageView) SetFocused(focused bool) {
	iv.focused = focused
}

// イメージ一覧を設定
func (iv *ImageView) SetImages(images []model.ImageInfo) {
	iv.images = images
	if iv.cursor >= len(images) {
		iv.cursor = len(images) - 1
	}
	if iv.cursor < 0 {
		iv.cursor = 0
	}
	iv.refreshContent()
}

// クリア
func (iv *ImageView) Clear() {
	iv.images = nil
	iv.cursor = 0
	iv.refreshContent()
}

// 選択中のイメージを返す
func (iv *ImageView) SelectedImage() *model.ImageInfo {
	if iv.cursor >= 0 && iv.cursor < len(iv.images) {
		img := iv.images[iv.cursor]
		return &img
	}
	return nil
}

// キー入力を処理
func (iv *ImageView) Update(msg tea.Msg) tea.Cmd {
	if !iv.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, iv.keymap.Up):
			if iv.cursor > 0 {
				iv.cursor--
				iv.refreshContent()
				iv.ensureVisible()
			}
			return nil
		case key.Matches(msg, iv.keymap.Down):
			if iv.cursor < len(iv.images)-1 {
				iv.cursor++
				iv.refreshContent()
				iv.ensureVisible()
			}
			return nil
		}
	}

	var cmd tea.Cmd
	iv.viewport, cmd = iv.viewport.Update(msg)
	return cmd
}

// スクロール位置を調整してカーソルが見えるようにする
func (iv *ImageView) ensureVisible() {
	vpHeight := iv.viewport.Height
	if vpHeight < 1 {
		return
	}
	if iv.cursor < iv.viewport.YOffset {
		iv.viewport.SetYOffset(iv.cursor)
	}
	if iv.cursor >= iv.viewport.YOffset+vpHeight {
		iv.viewport.SetYOffset(iv.cursor - vpHeight + 1)
	}
}

// イメージ一覧パネルを描画
func (iv ImageView) View() string {
	header := iv.styles.Subtitle.Render(i18n.T("images.title"))

	if len(iv.images) == 0 {
		return header + "\n" + iv.styles.Muted.Render(i18n.T("images.no_images"))
	}

	// テーブルヘッダー
	tableHeader := iv.styles.Bold.Render(iv.formatRow(
		i18n.T("images.repo"),
		i18n.T("images.size"),
		i18n.T("images.created"),
		i18n.T("images.status"),
	))

	return header + "\n" + tableHeader + "\n" + iv.viewport.View()
}

func (iv *ImageView) refreshContent() {
	if len(iv.images) == 0 {
		iv.viewport.SetContent("")
		return
	}

	var lines []string
	for idx, img := range iv.images {
		name := img.DisplayName()
		size := img.SizeHuman()
		created := formatTimeAgo(img.Created)

		var status string
		if img.IsUnused() {
			if len(img.RepoTags) == 0 || img.RepoTags[0] == "<none>:<none>" {
				status = iv.styles.Warning.Render(i18n.T("images.dangling"))
			} else {
				status = iv.styles.Muted.Render(i18n.T("images.unused"))
			}
		} else {
			status = iv.styles.Running.Render(i18n.T("images.used"))
		}

		row := iv.formatRow(name, size, created, status)

		if idx == iv.cursor && iv.focused {
			lines = append(lines, iv.styles.SelectedItem.Width(iv.width-1).Render(row))
		} else {
			lines = append(lines, iv.styles.NormalItem.Width(iv.width-1).Render(row))
		}
	}

	iv.viewport.SetContent(strings.Join(lines, "\n"))
}

func (iv ImageView) formatRow(col1, col2, col3, col4 string) string {
	// カラム幅を動的に計算
	w := iv.width
	if w < 40 {
		w = 40
	}
	sizeW := 10
	createdW := 14
	statusW := 10
	nameW := w - sizeW - createdW - statusW - 6 // マージン分
	if nameW < 10 {
		nameW = 10
	}

	return fmt.Sprintf("%-*s  %*s  %-*s  %s", nameW, col1, sizeW, col2, createdW, col3, col4)
}

// 時間を「X ago」形式で返す
func formatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}
