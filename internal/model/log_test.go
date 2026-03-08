package model

import (
	"testing"
	"time"
)

func TestRingBuffer_Add_and_Entries(t *testing.T) {
	rb := NewRingBuffer(3)

	// 空のバッファ
	if got := rb.Count(); got != 0 {
		t.Errorf("Count() = %d, want 0", got)
	}

	// エントリ追加
	for i := 0; i < 3; i++ {
		rb.Add(LogEntry{Message: string(rune('a' + i))})
	}
	if got := rb.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}

	entries := rb.Entries()
	if len(entries) != 3 {
		t.Fatalf("Entries() len = %d, want 3", len(entries))
	}
	for i, e := range entries {
		want := string(rune('a' + i))
		if e.Message != want {
			t.Errorf("entries[%d].Message = %q, want %q", i, e.Message, want)
		}
	}
}

func TestRingBuffer_Wrap(t *testing.T) {
	rb := NewRingBuffer(3)

	// 5エントリ追加（2周分超）
	for i := 0; i < 5; i++ {
		rb.Add(LogEntry{Message: string(rune('a' + i))})
	}

	if got := rb.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}

	entries := rb.Entries()
	// c, d, e が残るべき
	expected := []string{"c", "d", "e"}
	for i, e := range entries {
		if e.Message != expected[i] {
			t.Errorf("entries[%d].Message = %q, want %q", i, e.Message, expected[i])
		}
	}
}

func TestRingBuffer_Clear(t *testing.T) {
	rb := NewRingBuffer(5)
	rb.Add(LogEntry{Message: "test"})
	rb.Clear()

	if got := rb.Count(); got != 0 {
		t.Errorf("Count() after Clear() = %d, want 0", got)
	}
}

func TestLogFilter_Match(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Service:   "web",
		Stream:    StreamStdout,
		Message:   "GET /api/health 200",
	}

	tests := []struct {
		name   string
		filter LogFilter
		want   bool
	}{
		{
			name:   "デフォルトフィルタは全通過",
			filter: DefaultLogFilter(),
			want:   true,
		},
		{
			name:   "stderr非表示でstdoutは通過",
			filter: LogFilter{ShowStdout: true, ShowStderr: false},
			want:   true,
		},
		{
			name:   "stdout非表示でstdoutは不通過",
			filter: LogFilter{ShowStdout: false, ShowStderr: true},
			want:   false,
		},
		{
			name:   "検索一致",
			filter: LogFilter{ShowStdout: true, ShowStderr: true, SearchQuery: "health"},
			want:   true,
		},
		{
			name:   "検索不一致",
			filter: LogFilter{ShowStdout: true, ShowStderr: true, SearchQuery: "error"},
			want:   false,
		},
		{
			name:   "サービスフィルタ一致",
			filter: LogFilter{ShowStdout: true, ShowStderr: true, Services: []string{"web"}},
			want:   true,
		},
		{
			name:   "サービスフィルタ不一致",
			filter: LogFilter{ShowStdout: true, ShowStderr: true, Services: []string{"db"}},
			want:   false,
		},
		{
			name:   "大文字小文字を無視した検索",
			filter: LogFilter{ShowStdout: true, ShowStderr: true, SearchQuery: "HEALTH"},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Match(entry); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_FilteredEntries(t *testing.T) {
	rb := NewRingBuffer(100)
	rb.Add(LogEntry{Service: "web", Stream: StreamStdout, Message: "GET /health"})
	rb.Add(LogEntry{Service: "web", Stream: StreamStderr, Message: "error: timeout"})
	rb.Add(LogEntry{Service: "db", Stream: StreamStdout, Message: "query executed"})

	// stderrのみ
	filter := LogFilter{ShowStdout: false, ShowStderr: true}
	entries := rb.FilteredEntries(filter)
	if len(entries) != 1 {
		t.Fatalf("FilteredEntries() len = %d, want 1", len(entries))
	}
	if entries[0].Message != "error: timeout" {
		t.Errorf("got message %q", entries[0].Message)
	}

	// サービスフィルタ
	filter = LogFilter{ShowStdout: true, ShowStderr: true, Services: []string{"db"}}
	entries = rb.FilteredEntries(filter)
	if len(entries) != 1 {
		t.Fatalf("FilteredEntries() len = %d, want 1", len(entries))
	}
}
