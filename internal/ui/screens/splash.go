package screens

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/r7sqtr/orca/internal/docker"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/ui"
)

// 起動/接続確認画面
type Splash struct {
	styles    ui.Styles
	client    *docker.Client
	width     int
	height    int
	err       error
	diag      docker.Diagnosis // エラー診断結果のキャッシュ
	connected bool
	dots      int
}

// Splashを作成
func NewSplash(styles ui.Styles) Splash {
	return Splash{styles: styles}
}

// サイズを設定
func (s *Splash) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Docker接続を試みるコマンドを返す
func ConnectCmd(dockerHost string) tea.Cmd {
	return func() tea.Msg {
		client, err := docker.NewClient(dockerHost)
		if err != nil {
			return ui.DockerConnectionFailedMsg{Err: err}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx); err != nil {
			client.Close()
			return ui.DockerConnectionFailedMsg{Err: err}
		}

		return splashConnectedMsg{client: client}
	}
}

// Splash内部用の接続成功メッセージ
type splashConnectedMsg struct {
	client *docker.Client
}

// メッセージを処理
func (s *Splash) Update(msg tea.Msg) (*docker.Client, tea.Cmd) {
	switch msg := msg.(type) {
	case splashConnectedMsg:
		s.connected = true
		s.client = msg.client
		return msg.client, nil
	case ui.DockerConnectionFailedMsg:
		s.err = msg.Err
		s.diag = docker.DiagnoseConnectionError(msg.Err)
		return nil, nil
	case tickDots:
		s.dots = (s.dots + 1) % 4
		return nil, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
			return tickDots{}
		})
	}
	return nil, nil
}

type tickDots struct{}

// ドットアニメーションを開始
func StartAnimation() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickDots{}
	})
}

// スプラッシュ画面を描画
func (s Splash) View() string {
	var content string

	if s.err != nil {
		title := s.styles.Error.Render(i18n.T("app.no_docker"))

		cause := s.styles.Muted.Render(i18n.TF("diag.cause", i18n.T(s.diag.Cause)))

		hints := s.styles.Muted.Render(i18n.T("diag.hints"))
		for _, h := range s.diag.Hints {
			hints += "\n" + s.styles.Muted.Render(fmt.Sprintf("  • %s", i18n.T(h)))
		}

		errMsg := s.styles.Muted.Render(s.err.Error())
		content = title + "\n\n" + cause + "\n\n" + hints + "\n\n" + errMsg + "\n\n" + s.styles.Muted.Render("[q] " + i18n.T("help.quit"))
	} else {
		dots := ""
		for i := 0; i < s.dots; i++ {
			dots += "."
		}
		content = s.styles.Title.Render(i18n.T("app.title")) + "\n\n" +
			s.styles.Muted.Render(i18n.T("app.connecting")+dots)
	}

	return lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, content)
}

// 接続済みかを返す
func (s Splash) IsConnected() bool {
	return s.connected
}
