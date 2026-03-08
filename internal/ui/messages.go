package ui

import (
	"github.com/vvsaito/orca/internal/docker"
	"github.com/vvsaito/orca/internal/model"
)

// DockerConnectedMsg はDocker接続成功メッセージ
type DockerConnectedMsg struct{}

// DockerConnectionFailedMsg はDocker接続失敗メッセージ
type DockerConnectionFailedMsg struct {
	Err error
}

// ProjectsLoadedMsg はプロジェクト一覧取得完了メッセージ
type ProjectsLoadedMsg struct {
	Projects []model.ComposeProject
}

// ProjectsLoadFailedMsg はプロジェクト一覧取得失敗メッセージ
type ProjectsLoadFailedMsg struct {
	Err error
}

// LogEntryMsg はログエントリ受信メッセージ
type LogEntryMsg struct {
	Entry model.LogEntry
}

// DockerEventMsg はDockerイベント受信メッセージ
type DockerEventMsg struct {
	Event docker.DockerEvent
}

// ComposeActionDoneMsg はCompose操作完了メッセージ
type ComposeActionDoneMsg struct {
	Action  docker.ComposeAction
	Target  string
	Err     error
}

// ExecRequestMsg はシェル接続要求メッセージ
type ExecRequestMsg struct {
	WorkingDir string
	Service    string
}

// LogExportedMsg はログエクスポート完了メッセージ
type LogExportedMsg struct {
	Path string
	Err  error
}

// LogCopiedMsg はログコピー完了メッセージ
type LogCopiedMsg struct {
	Err error
}

// EnvVarsLoadedMsg は環境変数取得完了メッセージ
type EnvVarsLoadedMsg struct {
	Vars []string
	Err  error
}

// TickMsg は定期更新メッセージ
type TickMsg struct{}
