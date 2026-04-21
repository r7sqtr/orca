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
	Projects      []model.ComposeProject
	ResolvedPaths []docker.ResolveResult
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

// イメージ一覧取得完了メッセージ
type ImagesLoadedMsg struct {
	Images []model.ImageInfo
	Err    error
}

// ボリューム一覧取得完了メッセージ
type VolumesLoadedMsg struct {
	Volumes []model.VolumeInfo
	Err     error
}

// イメージ削除完了メッセージ
type ImageRemovedMsg struct {
	ImageID string
	Err     error
}

// ボリューム削除完了メッセージ
type VolumeRemovedMsg struct {
	VolumeName string
	Err        error
}

// イメージ一括削除完了メッセージ
type ImagesPrunedMsg struct {
	SpaceReclaimed uint64
	Err            error
}

// ボリューム一括削除完了メッセージ
type VolumesPrunedMsg struct {
	SpaceReclaimed uint64
	Err            error
}

// 定期更新メッセージ
type TickMsg struct{}
