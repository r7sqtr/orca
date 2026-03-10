package ui

// UIのレイアウト計算を保持
type Layout struct {
	Width  int
	Height int

	SidebarWidth int
	DetailWidth  int

	HeaderHeight    int
	StatusBarHeight int
	HelpBarHeight   int
	ContentHeight   int

	ShowSidebar bool
}

const (
	minSidebarWidth  = 20
	maxSidebarWidth  = 40
	minWidthSidebar  = 60 // サイドバー表示の最小幅
	headerHeight     = 0
	statusBarHeight  = 1
	helpBarHeight    = 1
)

// ターミナルサイズに基づきレイアウトを計算
func CalcLayout(width, height int) Layout {
	l := Layout{
		Width:           width,
		Height:          height,
		HeaderHeight:    headerHeight,
		StatusBarHeight: statusBarHeight,
		HelpBarHeight:   helpBarHeight,
	}

	// サイドバーの幅を計算
	l.ShowSidebar = width >= minWidthSidebar
	if l.ShowSidebar {
		// ターミナル幅の1/4 (clamp: 20〜40)
		l.SidebarWidth = width / 4
		if l.SidebarWidth < minSidebarWidth {
			l.SidebarWidth = minSidebarWidth
		}
		if l.SidebarWidth > maxSidebarWidth {
			l.SidebarWidth = maxSidebarWidth
		}
		// ボーダー分を含む
		l.DetailWidth = width - l.SidebarWidth - 1
	} else {
		l.SidebarWidth = 0
		l.DetailWidth = width
	}

	// コンテンツ領域の高さ
	l.ContentHeight = height - l.HeaderHeight - l.StatusBarHeight - l.HelpBarHeight
	if l.ContentHeight < 0 {
		l.ContentHeight = 0
	}

	return l
}
