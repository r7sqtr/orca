package docker

import "strings"

// エラー診断結果
type Diagnosis struct {
	Cause string   // 推定原因のi18nキー
	Hints []string // 確認事項のi18nキー一覧
}

// Docker接続エラーを診断
func DiagnoseConnectionError(err error) Diagnosis {
	if err == nil {
		return Diagnosis{}
	}
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "is the docker daemon running"):
		return Diagnosis{
			Cause: "diag.conn.cause.not_running",
			Hints: []string{
				"diag.conn.hint.start_docker",
			},
		}
	case strings.Contains(msg, "permission denied"):
		return Diagnosis{
			Cause: "diag.conn.cause.permission",
			Hints: []string{
				"diag.conn.hint.add_group",
				"diag.conn.hint.relogin",
			},
		}
	case strings.Contains(msg, "no such file or directory"):
		return Diagnosis{
			Cause: "diag.conn.cause.no_socket",
			Hints: []string{
				"diag.conn.hint.check_installed",
				"diag.conn.hint.start_docker",
			},
		}
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "context deadline exceeded"):
		return Diagnosis{
			Cause: "diag.conn.cause.timeout",
			Hints: []string{
				"diag.conn.hint.check_load",
				"diag.conn.hint.check_resources",
			},
		}
	case strings.Contains(msg, "client version") && strings.Contains(msg, "is too new"):
		return Diagnosis{
			Cause: "diag.conn.cause.version_mismatch",
			Hints: []string{
				"diag.conn.hint.update_docker",
			},
		}
	default:
		return Diagnosis{
			Cause: "diag.conn.cause.unknown",
			Hints: []string{
				"diag.conn.hint.run_docker_info",
			},
		}
	}
}

// Compose実行エラーを診断
func DiagnoseComposeError(err error) Diagnosis {
	if err == nil {
		return Diagnosis{}
	}
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "port is already allocated") || strings.Contains(msg, "address already in use"):
		return Diagnosis{
			Cause: "diag.compose.cause.port_conflict",
			Hints: []string{
				"diag.compose.hint.check_port",
			},
		}
	case strings.Contains(msg, "image not found") || strings.Contains(msg, "pull access denied"):
		return Diagnosis{
			Cause: "diag.compose.cause.image_not_found",
			Hints: []string{
				"diag.compose.hint.check_image",
				"diag.compose.hint.docker_login",
			},
		}
	case strings.Contains(msg, "no such file or directory") && (strings.Contains(msg, "dockerfile") || strings.Contains(msg, "build")):
		return Diagnosis{
			Cause: "diag.compose.cause.file_not_found",
			Hints: []string{
				"diag.compose.hint.check_path",
			},
		}
	case strings.Contains(msg, "network") && strings.Contains(msg, "not found"):
		return Diagnosis{
			Cause: "diag.compose.cause.network_not_found",
			Hints: []string{
				"diag.compose.hint.create_network",
			},
		}
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "context deadline exceeded"):
		return Diagnosis{
			Cause: "diag.compose.cause.timeout",
			Hints: []string{
				"diag.compose.hint.check_network",
				"diag.compose.hint.check_daemon",
			},
		}
	case strings.Contains(msg, "no configuration file"):
		return Diagnosis{
			Cause: "diag.compose.cause.no_config",
			Hints: []string{
				"diag.compose.hint.check_compose_file",
			},
		}
	default:
		// 原因不明の場合は空のDiagnosisを返す（生エラーをそのまま使う）
		return Diagnosis{}
	}
}
