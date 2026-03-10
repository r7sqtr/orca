package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r7sqtr/orca/internal/config"
	"github.com/r7sqtr/orca/internal/docker"
	"github.com/r7sqtr/orca/internal/model"
	"github.com/r7sqtr/orca/internal/ui"
	"github.com/r7sqtr/orca/internal/ui/screens"
)

// アプリケーションの画面状態
type AppState int

const (
	StateSplash AppState = iota
	StateMain
)

// アプリケーションのルートモデル
type AppModel struct {
	state  AppState
	styles ui.Styles
	keymap ui.KeyMap
	cfg    model.AppConfig
	width  int
	height int

	splash screens.Splash
	main   *screens.MainScreen
	client *docker.Client
}

// AppModelを作成
func NewAppModel(cfg model.AppConfig) AppModel {
	theme := config.GetTheme(cfg.Theme)
	styles := ui.NewStyles(theme)
	keymap := ui.DefaultKeyMap()

	return AppModel{
		state:  StateSplash,
		styles: styles,
		keymap: keymap,
		cfg:    cfg,
		splash: screens.NewSplash(styles),
	}
}

// アプリケーションを初期化
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		screens.ConnectCmd(),
		screens.StartAnimation(),
	)
}

// メッセージを処理
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.splash.SetSize(msg.Width, msg.Height)
		if m.main != nil {
			m.main.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.cleanup()
			return m, tea.Quit
		}
		if m.state == StateSplash && msg.String() == "q" {
			m.cleanup()
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateSplash:
		return m.updateSplash(msg)
	case StateMain:
		return m.updateMain(msg)
	}

	return m, nil
}

func (m AppModel) updateSplash(msg tea.Msg) (tea.Model, tea.Cmd) {
	client, cmd := m.splash.Update(msg)
	if client != nil {
		m.client = client
		m.state = StateMain
		main := screens.NewMainScreen(m.styles, m.keymap, client, m.cfg)
		main.SetSize(m.width, m.height)
		m.main = &main
		return m, m.main.Init()
	}
	return m, cmd
}

func (m AppModel) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.main == nil {
		return m, nil
	}
	cmd := m.main.Update(msg)
	return m, cmd
}

func (m *AppModel) cleanup() {
	if m.main != nil {
		m.main.Cleanup()
	}
	if m.client != nil {
		m.client.Close()
	}
}

// アプリケーションを描画
func (m AppModel) View() string {
	switch m.state {
	case StateSplash:
		return m.splash.View()
	case StateMain:
		if m.main != nil {
			return m.main.View()
		}
	}
	return ""
}
