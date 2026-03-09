package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/model"
	"github.com/vvsaito/orca/internal/ui"
)

// SidebarItem はサイドバーのアイテム
type SidebarItem struct {
	Type        SidebarItemType
	ProjectName string
	ServiceName string
	Container   *model.ContainerStatus
}

// SidebarItemType はアイテム種別
type SidebarItemType int

const (
	ItemProject SidebarItemType = iota
	ItemService
)

// Sidebar は左サイドバーパネル
type Sidebar struct {
	styles    ui.Styles
	keymap    ui.KeyMap
	projects  []model.ComposeProject
	items     []SidebarItem
	collapsed map[string]bool // プロジェクト名 → 折りたたみ状態
	cursor    int
	offset    int // スクロールオフセット
	width     int
	height    int
	focused   bool
}

// NewSidebar はSidebarを作成する
func NewSidebar(styles ui.Styles, keymap ui.KeyMap) Sidebar {
	return Sidebar{
		styles:    styles,
		keymap:    keymap,
		collapsed: make(map[string]bool),
	}
}

// SetSize はサイズを設定する
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// SetFocused はフォーカス状態を設定する
func (s *Sidebar) SetFocused(focused bool) {
	s.focused = focused
}

// IsFocused はフォーカス状態を返す
func (s Sidebar) IsFocused() bool {
	return s.focused
}

// SetProjects はプロジェクト一覧を設定する
func (s *Sidebar) SetProjects(projects []model.ComposeProject) {
	s.projects = projects
	s.rebuildItems()
}

// rebuildItems は折りたたみ状態を参照してアイテムリストを再構築する
func (s *Sidebar) rebuildItems() {
	s.items = nil
	for _, proj := range s.projects {
		s.items = append(s.items, SidebarItem{
			Type:        ItemProject,
			ProjectName: proj.Name,
		})
		if !s.collapsed[proj.Name] {
			for _, svc := range proj.Services {
				s.items = append(s.items, SidebarItem{
					Type:        ItemService,
					ProjectName: proj.Name,
					ServiceName: svc.Name,
					Container:   svc.Container,
				})
			}
		}
	}

	// カーソルの有効範囲を維持
	if s.cursor >= len(s.items) {
		s.cursor = len(s.items) - 1
	}
	if s.cursor < 0 {
		s.cursor = 0
	}
}

// ToggleCollapse は選択中プロジェクトの折りたたみ状態をトグルする
func (s *Sidebar) ToggleCollapse() {
	item := s.SelectedItem()
	if item == nil {
		return
	}

	projectName := item.ProjectName
	s.collapsed[projectName] = !s.collapsed[projectName]

	// 折りたたみ時にカーソルがサービス上なら、プロジェクト行に移動
	if s.collapsed[projectName] && item.Type == ItemService {
		for i, it := range s.items {
			if it.Type == ItemProject && it.ProjectName == projectName {
				s.cursor = i
				break
			}
		}
	}

	s.rebuildItems()
	s.ensureVisible()
}

// SelectedItem は選択中のアイテムを返す
func (s Sidebar) SelectedItem() *SidebarItem {
	if s.cursor >= 0 && s.cursor < len(s.items) {
		item := s.items[s.cursor]
		return &item
	}
	return nil
}

// Update はキー入力を処理する
func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	if !s.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keymap.Up):
			s.moveUp()
		case key.Matches(msg, s.keymap.Down):
			s.moveDown()
		}
	}
	return nil
}

func (s *Sidebar) moveUp() {
	if s.cursor > 0 {
		s.cursor--
		s.ensureVisible()
	}
}

func (s *Sidebar) moveDown() {
	if s.cursor < len(s.items)-1 {
		s.cursor++
		s.ensureVisible()
	}
}

func (s *Sidebar) ensureVisible() {
	listHeight := s.height - 1 // タイトル行分
	if listHeight < 1 {
		listHeight = 1
	}
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if s.cursor >= s.offset+listHeight {
		s.offset = s.cursor - listHeight + 1
	}
}

// View はサイドバーを描画する
func (s Sidebar) View() string {
	// タイトル行（フォーカス状態を表示）
	title := i18n.T("sidebar.title")
	var titleLine string
	if s.focused {
		titleLine = s.styles.Title.Render(title + " ◀")
	} else {
		titleLine = s.styles.Muted.Render(title)
	}

	if len(s.items) == 0 {
		msg := i18n.T("sidebar.no_projects")
		return titleLine + "\n" + s.styles.Muted.Width(s.width).Height(s.height-1).Render(msg)
	}

	var lines []string
	lines = append(lines, titleLine)

	listHeight := s.height - 1 // タイトル行分を引く
	end := s.offset + listHeight
	if end > len(s.items) {
		end = len(s.items)
	}

	for idx := s.offset; idx < end; idx++ {
		item := s.items[idx]
		selected := idx == s.cursor

		var line string
		switch item.Type {
		case ItemProject:
			icon := "▼"
			if s.collapsed[item.ProjectName] {
				icon = "▶"
			}
			line = fmt.Sprintf("%s %s", icon, item.ProjectName)
			if selected && s.focused {
				line = s.styles.SelectedItem.Width(s.width - 1).Render(line)
			} else {
				line = s.styles.ProjectItem.Width(s.width - 1).Render(line)
			}
		case ItemService:
			icon := s.serviceIcon(item.Container)
			line = fmt.Sprintf("  %s %s", icon, item.ServiceName)
			if selected && s.focused {
				line = s.styles.SelectedItem.Width(s.width - 1).Render(line)
			} else {
				line = s.styles.NormalItem.Width(s.width - 1).Render(line)
			}
		}

		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")

	// 高さが足りない場合はパディング（タイトル行含む全体高さ）
	rendered := lipgloss.Height(content)
	if rendered < s.height {
		content += strings.Repeat("\n", s.height-rendered)
	}

	return content
}

func (s Sidebar) serviceIcon(ctr *model.ContainerStatus) string {
	if ctr == nil {
		return s.styles.Muted.Render("○")
	}
	switch ctr.State {
	case "running":
		return s.styles.Running.Render("●")
	case "paused":
		return s.styles.Warning.Render("◉")
	case "restarting":
		return s.styles.Warning.Render("◌")
	default:
		return s.styles.Stopped.Render("○")
	}
}
