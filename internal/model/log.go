package model

import (
	"strings"
	"sync"
	"time"
)

// ログのストリーム種別
type StreamType int

const (
	StreamStdout StreamType = iota
	StreamStderr
)

// ログの1行を表す
type LogEntry struct {
	Timestamp   time.Time
	Service     string
	Stream      StreamType
	Message     string
	ContainerID string
}

// ログのフィルタ条件
type LogFilter struct {
	SearchQuery string
	ShowStdout  bool
	ShowStderr  bool
	Services    []string // 空なら全サービス
}

// デフォルトのフィルタを返す
func DefaultLogFilter() LogFilter {
	return LogFilter{
		ShowStdout: true,
		ShowStderr: true,
	}
}

// エントリがフィルタに一致するかを返す
func (f LogFilter) Match(entry LogEntry) bool {
	if !f.ShowStdout && entry.Stream == StreamStdout {
		return false
	}
	if !f.ShowStderr && entry.Stream == StreamStderr {
		return false
	}
	if len(f.Services) > 0 {
		found := false
		for _, s := range f.Services {
			if s == entry.Service {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if f.SearchQuery != "" {
		if !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(f.SearchQuery)) {
			return false
		}
	}
	return true
}

// 固定サイズのログバッファ
type RingBuffer struct {
	mu      sync.RWMutex
	entries []LogEntry
	size    int
	head    int
	count   int
}

// 指定サイズのRingBufferを作成
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		entries: make([]LogEntry, size),
		size:    size,
	}
}

// エントリを追加
func (rb *RingBuffer) Add(entry LogEntry) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.entries[rb.head] = entry
	rb.head = (rb.head + 1) % rb.size
	if rb.count < rb.size {
		rb.count++
	}
}

// 全エントリを古い順に返す
func (rb *RingBuffer) Entries() []LogEntry {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	result := make([]LogEntry, rb.count)
	if rb.count < rb.size {
		copy(result, rb.entries[:rb.count])
	} else {
		start := rb.head
		copy(result, rb.entries[start:])
		copy(result[rb.size-start:], rb.entries[:start])
	}
	return result
}

// フィルタに一致するエントリを返す
func (rb *RingBuffer) FilteredEntries(filter LogFilter) []LogEntry {
	entries := rb.Entries()
	result := make([]LogEntry, 0, len(entries))
	for _, e := range entries {
		if filter.Match(e) {
			result = append(result, e)
		}
	}
	return result
}

// バッファ内のエントリ数を返す
func (rb *RingBuffer) Count() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.count
}

// バッファをクリア
func (rb *RingBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.head = 0
	rb.count = 0
}
