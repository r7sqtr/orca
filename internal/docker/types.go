package docker

// Compose ラベル定数
const (
	LabelComposeProject    = "com.docker.compose.project"
	LabelComposeService    = "com.docker.compose.service"
	LabelComposeWorkingDir = "com.docker.compose.project.working_dir"
	LabelComposeConfigFile = "com.docker.compose.project.config_files"
)

// ComposeAction はComposeに対する操作種別
type ComposeAction string

const (
	ActionUp      ComposeAction = "up"
	ActionDown    ComposeAction = "down"
	ActionRestart ComposeAction = "restart"
	ActionBuild   ComposeAction = "build"
)
