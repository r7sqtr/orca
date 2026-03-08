package model

import (
	"strings"
	"sync"
	"time"
)

// StreamType はログのストリーム種別
type StreamType int

const (
	StreamStdout StreamType = iota
	StreamStderr
)

// LogEntry はログの1行を表す
type LogEntry struct {
	Timestamp   time.Time
	Service     string
	Stream      StreamType
	Message     string
	ContainerID string
}

// LogFilter はログのフィルタ条件
type LogFilter struct {
	SearchQuery string
	ShowStdout  bool
	ShowStderr  bool
	Services    []string // 空なら全サービス
}

// DefaultLogFilter はデフォルトのフィルタを返す
func DefaultLogFilter() LogFilter {
	return LogFilter{
		ShowStdout: true,
		ShowStderr: true,
	}
}

// Match はエントリがフィルタに一致するかを返す
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

// RingBuffer は固定サイズのログバッファ
type RingBuffer struct {
	mu      sync.RWMutex
	entries []LogEntry
	size    int
	head    int
	count   int
}

// NewRingBuffer は指定サイズのRingBufferを作成する
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		entries: make([]LogEntry, size),
		size:    size,
	}
}

// Add はエントリを追加する
func (rb *RingBuffer) Add(entry LogEntry) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.entries[rb.head] = entry
	rb.head = (rb.head + 1) % rb.size
	if rb.count < rb.size {
		rb.count++
	}
}

// Entries は全エントリを古い順に返す
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

// FilteredEntries はフィルタに一致するエントリを返す
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

// Count はバッファ内のエントリ数を返す
func (rb *RingBuffer) Count() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.count
}

// Clear はバッファをクリアする
func (rb *RingBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.head = 0
	rb.count = 0
}
