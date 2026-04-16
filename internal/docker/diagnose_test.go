package docker

import (
	"errors"
	"testing"
)

func TestDiagnoseConnectionError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantCause string
		wantHints int
	}{
		{
			name:      "nil error",
			err:       nil,
			wantCause: "",
			wantHints: 0,
		},
		{
			name:      "connection refused",
			err:       errors.New("dial tcp 127.0.0.1:2375: connect: connection refused"),
			wantCause: "diag.conn.cause.not_running",
			wantHints: 1,
		},
		{
			name:      "daemon not running",
			err:       errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"),
			wantCause: "diag.conn.cause.not_running",
			wantHints: 1,
		},
		{
			name:      "permission denied",
			err:       errors.New("Got permission denied while trying to connect to the Docker daemon socket"),
			wantCause: "diag.conn.cause.permission",
			wantHints: 2,
		},
		{
			name:      "socket not found",
			err:       errors.New("dial unix /var/run/docker.sock: connect: no such file or directory"),
			wantCause: "diag.conn.cause.no_socket",
			wantHints: 2,
		},
		{
			name:      "timeout",
			err:       errors.New("request timeout"),
			wantCause: "diag.conn.cause.timeout",
			wantHints: 2,
		},
		{
			name:      "context deadline exceeded",
			err:       errors.New("context deadline exceeded"),
			wantCause: "diag.conn.cause.timeout",
			wantHints: 2,
		},
		{
			name:      "version mismatch",
			err:       errors.New("Error response from daemon: client version 1.44 is too new"),
			wantCause: "diag.conn.cause.version_mismatch",
			wantHints: 1,
		},
		{
			name:      "unknown error",
			err:       errors.New("something unexpected happened"),
			wantCause: "diag.conn.cause.unknown",
			wantHints: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := DiagnoseConnectionError(tt.err)
			if diag.Cause != tt.wantCause {
				t.Errorf("Cause = %q, want %q", diag.Cause, tt.wantCause)
			}
			if len(diag.Hints) != tt.wantHints {
				t.Errorf("Hints count = %d, want %d (hints: %v)", len(diag.Hints), tt.wantHints, diag.Hints)
			}
		})
	}
}

func TestDiagnoseComposeError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantCause string
		wantHints int
	}{
		{
			name:      "nil error",
			err:       nil,
			wantCause: "",
			wantHints: 0,
		},
		{
			name:      "port already allocated",
			err:       errors.New("Bind for 0.0.0.0:8080 failed: port is already allocated"),
			wantCause: "diag.compose.cause.port_conflict",
			wantHints: 1,
		},
		{
			name:      "address already in use",
			err:       errors.New("listen tcp 0.0.0.0:3000: bind: address already in use"),
			wantCause: "diag.compose.cause.port_conflict",
			wantHints: 1,
		},
		{
			name:      "image not found",
			err:       errors.New("manifest for myapp:latest not found: manifest unknown: image not found"),
			wantCause: "diag.compose.cause.image_not_found",
			wantHints: 2,
		},
		{
			name:      "pull access denied",
			err:       errors.New("pull access denied for private-registry/myapp"),
			wantCause: "diag.compose.cause.image_not_found",
			wantHints: 2,
		},
		{
			name:      "dockerfile not found",
			err:       errors.New("failed to solve: failed to read dockerfile: open Dockerfile: no such file or directory"),
			wantCause: "diag.compose.cause.file_not_found",
			wantHints: 1,
		},
		{
			name:      "build context not found",
			err:       errors.New("failed to solve: build path ./app: no such file or directory"),
			wantCause: "diag.compose.cause.file_not_found",
			wantHints: 1,
		},
		{
			name:      "network not found",
			err:       errors.New("network my_network not found"),
			wantCause: "diag.compose.cause.network_not_found",
			wantHints: 1,
		},
		{
			name:      "timeout",
			err:       errors.New("context deadline exceeded while pulling image"),
			wantCause: "diag.compose.cause.timeout",
			wantHints: 2,
		},
		{
			name:      "no configuration file",
			err:       errors.New("no configuration file provided: not found"),
			wantCause: "diag.compose.cause.no_config",
			wantHints: 1,
		},
		{
			name:      "unknown error returns empty diagnosis",
			err:       errors.New("something unexpected"),
			wantCause: "",
			wantHints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := DiagnoseComposeError(tt.err)
			if diag.Cause != tt.wantCause {
				t.Errorf("Cause = %q, want %q", diag.Cause, tt.wantCause)
			}
			if len(diag.Hints) != tt.wantHints {
				t.Errorf("Hints count = %d, want %d (hints: %v)", len(diag.Hints), tt.wantHints, diag.Hints)
			}
		})
	}
}
