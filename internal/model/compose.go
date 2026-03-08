package model

// ComposeProject はDocker Composeプロジェクトを表す
type ComposeProject struct {
	Name       string
	WorkingDir string
	ConfigFile string
	Services   []Service
}

// Service はComposeサービスを表す
type Service struct {
	Name        string
	ProjectName string
	Container   *ContainerStatus
}

// IsRunning はサービスが実行中かどうかを返す
func (s Service) IsRunning() bool {
	return s.Container != nil && s.Container.State == "running"
}
