package model

import "time"

// Dockerボリュームの情報
type VolumeInfo struct {
	Name       string
	Driver     string
	MountPoint string
	CreatedAt  time.Time
	Labels     map[string]string
	RefCount   int64
}

// 未使用ボリュームかどうかを返す
func (v VolumeInfo) IsUnused() bool {
	return v.RefCount == 0
}
