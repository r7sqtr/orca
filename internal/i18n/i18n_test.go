package i18n

import "testing"

func TestT(t *testing.T) {
	SetLanguage("ja")
	got := T("app.title")
	want := "orca - Docker Compose マネージャー"
	if got != want {
		t.Errorf("T('app.title') = %q, want %q", got, want)
	}
}

func TestT_English(t *testing.T) {
	SetLanguage("en")
	defer SetLanguage("ja")

	got := T("app.title")
	want := "orca - Docker Compose Manager"
	if got != want {
		t.Errorf("T('app.title') = %q, want %q", got, want)
	}
}

func TestT_FallbackToKey(t *testing.T) {
	SetLanguage("ja")
	got := T("nonexistent.key")
	if got != "nonexistent.key" {
		t.Errorf("T('nonexistent.key') = %q, want 'nonexistent.key'", got)
	}
}

func TestTF(t *testing.T) {
	SetLanguage("ja")
	got := TF("confirm.up", "web")
	want := "web を起動しますか？"
	if got != want {
		t.Errorf("TF('confirm.up', 'web') = %q, want %q", got, want)
	}
}

func TestSetLanguage_Invalid(t *testing.T) {
	SetLanguage("ja")
	SetLanguage("invalid")
	// 無効な言語は無視され、jaのまま
	got := T("app.title")
	want := "orca - Docker Compose マネージャー"
	if got != want {
		t.Errorf("after invalid SetLanguage: T('app.title') = %q, want %q", got, want)
	}
}
