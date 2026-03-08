package model

import "time"

// ContainerStatus はコンテナの状態を表す
type ContainerStatus struct {
	ID          string
	Name        string
	Image       string
	State       string // running, exited, paused, etc.
	Status      string // "Up 2 hours", "Exited (0) 5 minutes ago"
	Health      string // healthy, unhealthy, starting, none
	Ports       []PortMapping
	CreatedAt   time.Time
	StartedAt   time.Time
	Labels      map[string]string
	ServiceName string
	ProjectName string
}

// PortMapping はポートマッピングを表す
type PortMapping struct {
	HostIP        string
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

// IsRunning はコンテナが実行中かどうかを返す
func (c ContainerStatus) IsRunning() bool {
	return c.State == "running"
}
