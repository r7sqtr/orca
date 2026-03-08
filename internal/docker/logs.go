package docker

import (
	"context"
	"encoding/binary"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/vvsaito/orca/internal/model"
)

// LogStreamer はコンテナログをストリーミングする
type LogStreamer struct {
	client *Client
	cancel context.CancelFunc
}

// NewLogStreamer はLogStreamerを作成する
func NewLogStreamer(client *Client) *LogStreamer {
	return &LogStreamer{client: client}
}

// Stream はコンテナログのストリーミングを開始する
// チャネルにLogEntryを送信する。コンテキストがキャンセルされると停止する
func (ls *LogStreamer) Stream(ctx context.Context, containerID, serviceName string, ch chan<- model.LogEntry) error {
	ctx, cancel := context.WithCancel(ctx)
	ls.cancel = cancel

	reader, err := ls.client.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
		Tail:       "500",
	})
	if err != nil {
		cancel()
		return err
	}

	go func() {
		defer reader.Close()
		defer cancel()
		ls.parseMultiplexedStream(reader, serviceName, containerID, ch)
	}()

	return nil
}

// Stop はストリーミングを停止する
func (ls *LogStreamer) Stop() {
	if ls.cancel != nil {
		ls.cancel()
	}
}

// parseMultiplexedStream はDockerのmultiplexed streamをパースする
// ヘッダフォーマット: [stream_type(1), 0, 0, 0, size(4)]
func (ls *LogStreamer) parseMultiplexedStream(reader io.Reader, serviceName, containerID string, ch chan<- model.LogEntry) {
	header := make([]byte, 8)

	for {
		_, err := io.ReadFull(reader, header)
		if err != nil {
			return
		}

		streamType := model.StreamStdout
		if header[0] == 2 {
			streamType = model.StreamStderr
		}

		size := binary.BigEndian.Uint32(header[4:8])
		if size == 0 {
			continue
		}

		payload := make([]byte, size)
		_, err = io.ReadFull(reader, payload)
		if err != nil {
			return
		}

		line := strings.TrimRight(string(payload), "\n\r")
		if line == "" {
			continue
		}

		entry := model.LogEntry{
			Service:     serviceName,
			Stream:      streamType,
			ContainerID: containerID,
		}

		// タイムスタンプのパース ("2024-01-01T00:00:00.000000000Z メッセージ")
		if idx := strings.IndexByte(line, ' '); idx > 0 {
			if t, err := time.Parse(time.RFC3339Nano, line[:idx]); err == nil {
				entry.Timestamp = t
				entry.Message = line[idx+1:]
			} else {
				entry.Timestamp = time.Now()
				entry.Message = line
			}
		} else {
			entry.Timestamp = time.Now()
			entry.Message = line
		}

		select {
		case ch <- entry:
		default:
			// チャネルがフルの場合はドロップ
		}
	}
}
