package model

import (
	"fmt"
	"time"
)

// Dockerイメージの情報
type ImageInfo struct {
	ID         string
	RepoTags   []string
	Size       int64
	Created    time.Time
	Containers int64
}

// 未使用イメージかどうかを返す
func (img ImageInfo) IsUnused() bool {
	return img.Containers == 0
}

// サイズを人間が読みやすい形式で返す
func (img ImageInfo) SizeHuman() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case img.Size >= GB:
		return fmt.Sprintf("%.1f GB", float64(img.Size)/float64(GB))
	case img.Size >= MB:
		return fmt.Sprintf("%.1f MB", float64(img.Size)/float64(MB))
	case img.Size >= KB:
		return fmt.Sprintf("%.1f KB", float64(img.Size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", img.Size)
	}
}

// 表示用のリポジトリ:タグを返す（dangling対応）
func (img ImageInfo) DisplayName() string {
	if len(img.RepoTags) == 0 || img.RepoTags[0] == "<none>:<none>" {
		return img.ID[:12]
	}
	return img.RepoTags[0]
}
