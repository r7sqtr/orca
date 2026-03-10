package ui

import (
	"github.com/r7sqtr/orca/internal/docker"
	"github.com/r7sqtr/orca/internal/model"
)

// Docker接続成功メッセージ
type DockerConnectedMsg struct{}

// Docker接続失敗メッセージ
type DockerConnectionFailedMsg struct {
	Err error
}

// プロジェクト一覧取得完了メッセージ
type ProjectsLoadedMsg struct {
	Projects []model.ComposeProject
}

// プロジェクト一覧取得失敗メッセージ
type ProjectsLoadFailedMsg struct {
	Err error
}

// ログエントリ受信メッセージ
type LogEntryMsg struct {
	Entry model.LogEntry
}

// Dockerイベント受信メッセージ
type DockerEventMsg struct {
	Event docker.DockerEvent
}

// Compose操作完了メッセージ
type ComposeActionDoneMsg struct {
	Action  docker.ComposeAction
	Target  string
	Err     error
}

// シェル接続要求メッセージ
type ExecRequestMsg struct {
	WorkingDir string
	Service    string
}

// ログエクスポート完了メッセージ
type LogExportedMsg struct {
	Path string
	Err  error
}

// ログコピー完了メッセージ
type LogCopiedMsg struct {
	Err error
}

// 環境変数取得完了メッセージ
type EnvVarsLoadedMsg struct {
	Vars []string
	Err  error
}

// 定期更新メッセージ
type TickMsg struct{}
