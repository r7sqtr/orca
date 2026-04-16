package docker

// Compose ラベル定数
const (
	LabelComposeProject    = "com.docker.compose.project"
	LabelComposeService    = "com.docker.compose.service"
	LabelComposeWorkingDir = "com.docker.compose.project.working_dir"
	LabelComposeConfigFile = "com.docker.compose.project.config_files"
)

// Composeに対する操作種別
type ComposeAction string

const (
	ActionUp      ComposeAction = "up"
	ActionStop    ComposeAction = "stop"
	ActionRestart ComposeAction = "restart"
	ActionBuild   ComposeAction = "build"
	ActionExec    ComposeAction = "exec"
	ActionDown    ComposeAction = "down"
)
