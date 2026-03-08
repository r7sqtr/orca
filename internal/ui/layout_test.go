package ui

import "testing"

func TestCalcLayout(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		wantSidebar   bool
	}{
		{"広いターミナル", 120, 40, true},
		{"狭いターミナル", 50, 30, false},
		{"最小サイドバー幅", 60, 30, true},
		{"ギリギリ非表示", 59, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := CalcLayout(tt.width, tt.height)
			if l.ShowSidebar != tt.wantSidebar {
				t.Errorf("ShowSidebar = %v, want %v", l.ShowSidebar, tt.wantSidebar)
			}
			if l.ShowSidebar {
				if l.SidebarWidth < minSidebarWidth {
					t.Errorf("SidebarWidth = %d, want >= %d", l.SidebarWidth, minSidebarWidth)
				}
				if l.SidebarWidth > maxSidebarWidth {
					t.Errorf("SidebarWidth = %d, want <= %d", l.SidebarWidth, maxSidebarWidth)
				}
				if l.DetailWidth <= 0 {
					t.Errorf("DetailWidth = %d, want > 0", l.DetailWidth)
				}
			}
			if l.ContentHeight <= 0 {
				t.Errorf("ContentHeight = %d, want > 0", l.ContentHeight)
			}
		})
	}
}
