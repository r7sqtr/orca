package screens

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vvsaito/orca/internal/docker"
	"github.com/vvsaito/orca/internal/i18n"
	"github.com/vvsaito/orca/internal/model"
	"github.com/vvsaito/orca/internal/ui"
	"github.com/vvsaito/orca/internal/ui/components"
	"github.com/vvsaito/orca/internal/ui/panels"
)

// FocusedPanel はフォーカス中のパネル
type FocusedPanel int

const (
	FocusSidebar FocusedPanel = iota
	FocusDetail
)

// MainScreen はメイン画面
type MainScreen struct {
	styles      ui.Styles
	keymap      ui.KeyMap
	layout      ui.Layout
	client      *docker.Client
	compose     *docker.ComposeExec
	sidebar     panels.Sidebar
	detail      panels.Detail
	statusBar   components.StatusBar
	helpBar     components.HelpBar
	confirm     components.ConfirmDialog
	helpOverlay components.HelpOverlay
	projects    []model.ComposeProject
	focused     FocusedPanel
	logCh       chan model.LogEntry
	eventCh     chan docker.DockerEvent
	streamer    *docker.LogStreamer
	cancelCtx   context.CancelFunc

	// 現在ログストリーミング中のコンテナID
	activeStreamContainerID string
	// 現在環境変数を読み込み済みのコンテナID
	activeEnvContainerID string
	// ヘルプオーバーレイ表示中フラグ
	showHelp bool
	// アプリケーション設定
	cfg model.AppConfig
}

// NewMainScreen はMainScreenを作成する
func NewMainScreen(styles ui.Styles, keymap ui.KeyMap, client *docker.Client, cfg model.AppConfig) MainScreen {
	ms := MainScreen{
		styles:      styles,
		keymap:      keymap,
		client:      client,
		compose:     docker.NewComposeExec(),
		sidebar:     panels.NewSidebar(styles, keymap),
		detail:      panels.NewDetail(styles, keymap, cfg.LogBufferSize),
		statusBar:   components.NewStatusBar(styles),
		helpBar:     components.NewHelpBar(styles),
		confirm:     components.NewConfirmDialog(styles),
		helpOverlay: components.NewHelpOverlay(styles),
		focused:     FocusSidebar,
		logCh:       make(chan model.LogEntry, 256),
		eventCh:     make(chan docker.DockerEvent, 64),
	}

	ms.cfg = cfg
	ms.sidebar.SetFocused(true)
	ms.statusBar.SetConnected(true)

	return ms
}

// Init はメイン画面を初期化する
func (ms *MainScreen) Init() tea.Cmd {
	return tea.Batch(
		ms.loadProjects(),
		ms.watchEvents(),
		ms.listenLogEntries(),
		ms.listenDockerEvents(),
		ms.tick(),
	)
}

// SetSize はサイズを設定する
func (ms *MainScreen) SetSize(width, height int) {
	ms.layout = ui.CalcLayout(width, height)
	ms.sidebar.SetSize(ms.layout.SidebarWidth, ms.layout.ContentHeight)
	ms.detail.SetSize(ms.layout.DetailWidth-1, ms.layout.ContentHeight) // PaddingLeft(1)分を引く
	ms.statusBar.SetSize(width)
	ms.helpBar.SetSize(width)
}

// Update はメッセージを処理する
func (ms *MainScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ms.handleKey(msg)
	case ui.ProjectsLoadedMsg:
		ms.projects = msg.Projects
		ms.sidebar.SetProjects(msg.Projects)
		return ms.updateSelectedService()
	case ui.LogEntryMsg:
		ms.detail.LogView().AddEntry(msg.Entry)
		return ms.listenLogEntries()
	case ui.DockerEventMsg:
		return tea.Batch(
			ms.loadProjects(),
			ms.listenDockerEvents(),
		)
	case ui.ComposeActionDoneMsg:
		if msg.Err != nil {
			ms.statusBar.SetMessage(i18n.TF("error.compose_exec", msg.Err.Error()))
		} else {
			ms.statusBar.SetMessage("")
		}
		return ms.loadProjects()
	case ui.LogCopiedMsg:
		if msg.Err != nil {
			ms.statusBar.SetMessage(i18n.TF("log.copy_failed", msg.Err.Error()))
		} else {
			ms.statusBar.SetMessage(i18n.T("log.copied"))
		}
		return nil
	case ui.LogExportedMsg:
		if msg.Err != nil {
			ms.statusBar.SetMessage(i18n.TF("log.export_failed", msg.Err.Error()))
		} else {
			ms.statusBar.SetMessage(i18n.TF("log.exported", msg.Path))
		}
		return nil
	case ui.EnvVarsLoadedMsg:
		if msg.Err != nil {
			ms.detail.EnvView().Clear()
		} else {
			ms.detail.EnvView().SetEnvVars(msg.Vars)
		}
		return nil
	case ui.TickMsg:
		return ms.tick()
	}

	// 非KeyMsgを子コンポーネントにも委譲
	switch ms.focused {
	case FocusDetail:
		return ms.detail.Update(msg)
	}

	return nil
}

func (ms *MainScreen) handleKey(msg tea.KeyMsg) tea.Cmd {
	// ヘルプオーバーレイがアクティブな場合
	if ms.showHelp {
		switch msg.String() {
		case "esc", "?", "q":
			ms.showHelp = false
			ms.helpOverlay.Hide()
		}
		return nil
	}

	// 確認ダイアログがアクティブな場合: y/n/esc のみ受付
	if ms.confirm.IsActive() {
		result := ms.confirm.Update(msg)
		if result != nil {
			if result.Confirmed {
				return ms.executeAction(result.Action, result.Target)
			}
			// キャンセル時もヘルプバーを通常モードに戻す
			ms.helpBar.SetMode(components.HelpModeNormal)
			ms.updateHelpMode()
		}
		return nil
	}

	// 検索バーがアクティブな場合: 全キーを検索バーに委譲
	if ms.detail.LogView().IsSearchActive() {
		return ms.detail.Update(msg)
	}

	// ctrl+c は常に終了
	if msg.String() == "ctrl+c" {
		return tea.Quit
	}

	// ヘルプ（?）は常に有効
	if key.Matches(msg, ms.keymap.Help) || msg.String() == "?" {
		ms.showHelp = true
		ms.helpOverlay.Show()
		return nil
	}

	// パネル切替（Ctrl+H/L）は常に有効
	if key.Matches(msg, ms.keymap.FocusLeft) {
		ms.setFocus(FocusSidebar)
		return nil
	}
	if key.Matches(msg, ms.keymap.FocusRight) {
		ms.setFocus(FocusDetail)
		return nil
	}

	// フォーカスに応じてキーを分岐
	switch ms.focused {
	case FocusSidebar:
		return ms.handleSidebarKey(msg)
	case FocusDetail:
		return ms.handleDetailKey(msg)
	}

	return nil
}

// handleSidebarKey はサイドバーフォーカス時のキー処理
func (ms *MainScreen) handleSidebarKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, ms.keymap.Quit):
		return tea.Quit
	case key.Matches(msg, ms.keymap.Start):
		return ms.confirmAction(docker.ActionUp)
	case key.Matches(msg, ms.keymap.Stop):
		return ms.confirmAction(docker.ActionDown)
	case key.Matches(msg, ms.keymap.Restart):
		return ms.confirmAction(docker.ActionRestart)
	case key.Matches(msg, ms.keymap.Build):
		return ms.confirmAction(docker.ActionBuild)
	case key.Matches(msg, ms.keymap.Exec):
		return ms.confirmAction(docker.ActionExec)
	case key.Matches(msg, ms.keymap.Info):
		return ms.switchDetailTab(panels.TabInfo)
	case key.Matches(msg, ms.keymap.Logs):
		return ms.switchDetailTab(panels.TabLogs)
	case key.Matches(msg, ms.keymap.EnvVars):
		return ms.switchDetailTab(panels.TabEnv)
	}

	cmd := ms.sidebar.Update(msg)
	serviceCmd := ms.updateSelectedService()
	if serviceCmd != nil {
		return tea.Batch(cmd, serviceCmd)
	}
	return cmd
}

// handleDetailKey はDetailパネルフォーカス時のキー処理
func (ms *MainScreen) handleDetailKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, ms.keymap.Back), key.Matches(msg, ms.keymap.Quit):
		ms.setFocus(FocusSidebar)
		return nil
	case key.Matches(msg, ms.keymap.Copy):
		if ms.detail.ActiveTab() == panels.TabLogs {
			return ms.copyLogs()
		}
	case key.Matches(msg, ms.keymap.Export):
		if ms.detail.ActiveTab() == panels.TabLogs {
			return ms.exportLogs()
		}
	case key.Matches(msg, ms.keymap.Info):
		return ms.switchDetailTab(panels.TabInfo)
	case key.Matches(msg, ms.keymap.Logs):
		return ms.switchDetailTab(panels.TabLogs)
	case key.Matches(msg, ms.keymap.EnvVars):
		return ms.switchDetailTab(panels.TabEnv)
	case key.Matches(msg, ms.keymap.Tab):
		return ms.cycleDetailTab()
	}

	return ms.detail.Update(msg)
}

// cycleDetailTab はDetailパネル内のタブを巡回する
func (ms *MainScreen) cycleDetailTab() tea.Cmd {
	next := (ms.detail.ActiveTab() + 1) % panels.NumDetailTabs
	return ms.switchDetailTab(next)
}

// setFocus はフォーカスを指定パネルに切り替える
func (ms *MainScreen) setFocus(panel FocusedPanel) {
	ms.focused = panel
	ms.sidebar.SetFocused(panel == FocusSidebar)
	ms.detail.SetFocused(panel == FocusDetail)
	ms.updateHelpMode()
}

// switchDetailTab はDetailタブを切り替え、必要に応じて環境変数をロードする
func (ms *MainScreen) switchDetailTab(tab panels.DetailTab) tea.Cmd {
	ms.detail.SwitchTab(tab)
	ms.updateHelpMode()
	if tab == panels.TabEnv {
		return ms.loadEnvVarsCmd()
	}
	return nil
}

func (ms *MainScreen) updateHelpMode() {
	switch ms.focused {
	case FocusDetail:
		switch ms.detail.ActiveTab() {
		case panels.TabLogs:
			ms.helpBar.SetMode(components.HelpModeLogs)
		case panels.TabEnv:
			ms.helpBar.SetMode(components.HelpModeEnv)
		default:
			ms.helpBar.SetMode(components.HelpModeInfo)
		}
	default:
		ms.helpBar.SetMode(components.HelpModeNormal)
	}
}

func (ms *MainScreen) updateSelectedService() tea.Cmd {
	item := ms.sidebar.SelectedItem()
	if item == nil {
		return nil
	}

	var cmds []tea.Cmd

	switch item.Type {
	case panels.ItemService:
		ms.detail.SetService(item.ProjectName, item.ServiceName, item.Container)
		ms.statusBar.SetContext(item.ProjectName, item.ServiceName)

		if item.Container != nil && item.Container.IsRunning() {
			if item.Container.ID != ms.activeStreamContainerID {
				ms.startLogStream(item.Container.ID, item.ServiceName)
			}
		} else if ms.activeStreamContainerID != "" {
			ms.stopLogStream()
		}

		// サービスが切り替わったら環境変数キャッシュをリセットし再取得
		if item.Container == nil || item.Container.ID != ms.activeEnvContainerID {
			ms.activeEnvContainerID = ""
			ms.detail.EnvView().Clear()
			if ms.detail.ActiveTab() == panels.TabEnv {
				cmds = append(cmds, ms.loadEnvVarsCmd())
			}
		}
	case panels.ItemProject:
		ms.statusBar.SetContext(item.ProjectName, "")
	}

	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

func (ms *MainScreen) startLogStream(containerID, serviceName string) {
	ms.stopLogStream()

	ctx, cancel := context.WithCancel(context.Background())
	ms.cancelCtx = cancel
	ms.activeStreamContainerID = containerID

	ms.detail.LogView().Clear()
	ms.streamer = docker.NewLogStreamer(ms.client)
	_ = ms.streamer.Stream(ctx, containerID, serviceName, ms.logCh)
}

func (ms *MainScreen) stopLogStream() {
	if ms.streamer != nil {
		ms.streamer.Stop()
		ms.streamer = nil
	}
	if ms.cancelCtx != nil {
		ms.cancelCtx()
		ms.cancelCtx = nil
	}
	ms.activeStreamContainerID = ""
}

// needsConfirm はアクションに確認ダイアログが必要か判定する
func (ms *MainScreen) needsConfirm(action docker.ComposeAction) bool {
	switch action {
	case docker.ActionUp:
		return ms.cfg.ConfirmActions.Up
	case docker.ActionDown:
		return ms.cfg.ConfirmActions.Down
	case docker.ActionRestart:
		return ms.cfg.ConfirmActions.Restart
	case docker.ActionBuild:
		return ms.cfg.ConfirmActions.Build
	case docker.ActionExec:
		return ms.cfg.ConfirmActions.Exec
	}
	return true
}

func (ms *MainScreen) confirmAction(action docker.ComposeAction) tea.Cmd {
	item := ms.sidebar.SelectedItem()
	if item == nil {
		return nil
	}

	// exec はサービスが起動中でないと実行不可
	if action == docker.ActionExec {
		if item.Type != panels.ItemService {
			return nil
		}
		if item.Container == nil || !item.Container.IsRunning() {
			ms.statusBar.SetMessage(i18n.T("exec.not_running"))
			return nil
		}
	}

	target := item.ProjectName
	if item.Type == panels.ItemService {
		target = item.ServiceName
	}

	// 確認不要の場合は直接実行
	if !ms.needsConfirm(action) {
		return ms.executeAction(string(action), target)
	}

	var msgKey string
	switch action {
	case docker.ActionUp:
		msgKey = "confirm.up"
	case docker.ActionDown:
		msgKey = "confirm.down"
	case docker.ActionRestart:
		msgKey = "confirm.restart"
	case docker.ActionBuild:
		msgKey = "confirm.build"
	case docker.ActionExec:
		msgKey = "confirm.exec"
	}

	ms.confirm.Show(
		i18n.TF(msgKey, target),
		string(action),
		target,
	)
	ms.helpBar.SetMode(components.HelpModeConfirm)
	return nil
}

func (ms *MainScreen) executeAction(action, target string) tea.Cmd {
	ms.helpBar.SetMode(components.HelpModeNormal)

	// exec は専用メソッドに委譲
	if docker.ComposeAction(action) == docker.ActionExec {
		return ms.execShell()
	}

	item := ms.sidebar.SelectedItem()
	if item == nil {
		return nil
	}

	workingDir := ms.findWorkingDir(item.ProjectName)
	composeAction := docker.ComposeAction(action)
	service := ""
	if item.Type == panels.ItemService {
		service = item.ServiceName
	}

	// ビルドはタイムアウトを長めに
	timeout := 30 * time.Second
	if composeAction == docker.ActionBuild {
		timeout = 5 * time.Minute
	}

	ms.statusBar.SetMessage(i18n.T("action."+action) + "...")

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := ms.compose.Run(ctx, workingDir, composeAction, service)
		return ui.ComposeActionDoneMsg{
			Action: composeAction,
			Target: target,
			Err:    err,
		}
	}
}

// execShell はシェル接続を実行する
func (ms *MainScreen) execShell() tea.Cmd {
	item := ms.sidebar.SelectedItem()
	if item == nil || item.Type != panels.ItemService {
		return nil
	}

	if item.Container == nil || !item.Container.IsRunning() {
		ms.statusBar.SetMessage(i18n.T("exec.not_running"))
		return nil
	}

	workingDir := ms.findWorkingDir(item.ProjectName)
	cmd := ms.compose.ExecCommand(workingDir, item.ServiceName)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ui.ComposeActionDoneMsg{
				Action: "exec",
				Target: item.ServiceName,
				Err:    err,
			}
		}
		return ui.ComposeActionDoneMsg{
			Action: "exec",
			Target: item.ServiceName,
		}
	})
}

// copyLogs はログをクリップボードにコピーする
func (ms *MainScreen) copyLogs() tea.Cmd {
	text := ms.detail.LogView().GetPlainText()
	if text == "" {
		return nil
	}

	return func() tea.Msg {
		err := clipboard.WriteAll(text)
		return ui.LogCopiedMsg{Err: err}
	}
}

// exportLogs はログをファイルにエクスポートする
func (ms *MainScreen) exportLogs() tea.Cmd {
	text := ms.detail.LogView().GetPlainText()
	if text == "" {
		return nil
	}

	item := ms.sidebar.SelectedItem()
	serviceName := "unknown"
	if item != nil && item.ServiceName != "" {
		serviceName = item.ServiceName
	}

	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			return ui.LogExportedMsg{Err: err}
		}

		dir := filepath.Join(home, ".config", "orca", "logs")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return ui.LogExportedMsg{Err: err}
		}

		filename := fmt.Sprintf("%s_%s.log", serviceName, time.Now().Format("20060102_150405"))
		path := filepath.Join(dir, filename)

		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			return ui.LogExportedMsg{Err: err}
		}

		return ui.LogExportedMsg{Path: path}
	}
}

// loadEnvVarsCmd は選択中サービスの環境変数をロードするCmdを返す
func (ms *MainScreen) loadEnvVarsCmd() tea.Cmd {
	item := ms.sidebar.SelectedItem()
	if item == nil || item.Type != panels.ItemService {
		ms.detail.EnvView().Clear()
		return nil
	}

	if item.Container == nil || !item.Container.IsRunning() {
		ms.detail.EnvView().SetEnvVars(nil)
		return nil
	}

	// 既に同じコンテナの環境変数を読み込み済みならスキップ
	if item.Container.ID == ms.activeEnvContainerID {
		return nil
	}
	ms.activeEnvContainerID = item.Container.ID

	containerID := item.Container.ID
	client := ms.client

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		vars, err := client.GetContainerEnv(ctx, containerID)
		return ui.EnvVarsLoadedMsg{Vars: vars, Err: err}
	}
}

func (ms *MainScreen) findWorkingDir(projectName string) string {
	for _, p := range ms.projects {
		if p.Name == projectName {
			return p.WorkingDir
		}
	}
	return ""
}

func (ms *MainScreen) loadProjects() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		containers, err := ms.client.ListComposeContainers(ctx)
		if err != nil {
			return ui.ProjectsLoadFailedMsg{Err: err}
		}

		projects := docker.GroupByProject(containers)
		return ui.ProjectsLoadedMsg{Projects: projects}
	}
}

func (ms *MainScreen) watchEvents() tea.Cmd {
	return func() tea.Msg {
		go docker.WatchEvents(context.Background(), ms.client, ms.eventCh)
		return nil
	}
}

func (ms *MainScreen) listenLogEntries() tea.Cmd {
	ch := ms.logCh
	return func() tea.Msg {
		entry := <-ch
		return ui.LogEntryMsg{Entry: entry}
	}
}

func (ms *MainScreen) listenDockerEvents() tea.Cmd {
	ch := ms.eventCh
	return func() tea.Msg {
		event := <-ch
		return ui.DockerEventMsg{Event: event}
	}
}

func (ms *MainScreen) tick() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return ui.TickMsg{}
	})
}

// Cleanup はリソースを解放する
func (ms *MainScreen) Cleanup() {
	ms.stopLogStream()
}

// View はメイン画面を描画する
func (ms MainScreen) View() string {
	// ヘルプオーバーレイ
	if ms.showHelp {
		return ms.helpOverlay.View(ms.layout.Width, ms.layout.Height)
	}

	// 確認ダイアログ（オーバーレイ）
	if ms.confirm.IsActive() {
		return ms.confirm.View(ms.layout.Width, ms.layout.Height)
	}

	w := ms.layout.Width
	ch := ms.layout.ContentHeight

	// コンテンツ領域
	var content string
	if ms.layout.ShowSidebar {
		sidebarView := ms.sidebar.View()
		detailView := ms.detail.View()

		// サイドバー: ボーダー色をフォーカスに応じて変更
		sidebarBorderColor := ms.styles.Theme.Border
		if ms.focused == FocusSidebar {
			sidebarBorderColor = ms.styles.Theme.Primary
		}

		// サイドバー（右ボーダー付き）
		sidebarRendered := lipgloss.NewStyle().
			Width(ms.layout.SidebarWidth).
			Height(ch).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(sidebarBorderColor).
			Render(sidebarView)

		// Detail（左パディングのみ）
		detailRendered := lipgloss.NewStyle().
			Width(ms.layout.DetailWidth).
			Height(ch).
			PaddingLeft(1).
			Render(detailView)

		content = lipgloss.JoinHorizontal(lipgloss.Top, sidebarRendered, detailRendered)
	} else {
		content = lipgloss.NewStyle().
			Width(w).
			Height(ch).
			PaddingLeft(1).
			Render(ms.detail.View())
	}

	// ステータスバー（1行固定）
	statusBar := ms.statusBar.View()

	// ヘルプバー（1行固定）
	helpBar := ms.helpBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, content, statusBar, helpBar)
}

// truncate は文字列を指定幅で切り詰める
func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	w := 0
	for i, r := range s {
		rw := 1
		if r > 0x7F {
			rw = 2 // CJK幅
		}
		if w+rw > maxWidth {
			return s[:i]
		}
		w += rw
	}
	return s
}

// pad は文字列を指定行数になるよう空行で埋める
func pad(s string, targetHeight int) string {
	lines := strings.Count(s, "\n") + 1
	if lines < targetHeight {
		s += strings.Repeat("\n", targetHeight-lines)
	}
	return s
}
