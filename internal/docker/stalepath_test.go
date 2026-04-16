package docker

import (
	"os"
	"path/filepath"
	"testing"
)

// ディレクトリ移動シナリオでfindRelocatedProjectが新パスを発見できること
func TestFindRelocatedProject_MovedToSubdir(t *testing.T) {
	// テスト用ディレクトリ構造:
	// tmpdir/
	//   Development/
	//     yy/
	//       myproject/
	//         docker-compose.yml
	tmpDir := t.TempDir()
	devDir := filepath.Join(tmpDir, "Development")
	newProjectDir := filepath.Join(devDir, "yy", "myproject")
	if err := os.MkdirAll(newProjectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	composeFile := filepath.Join(newProjectDir, "docker-compose.yml")
	if err := os.WriteFile(composeFile, []byte("version: '3'"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 旧パス: tmpdir/Development/myproject（存在しない）
	oldPath := filepath.Join(devDir, "myproject")
	info := ProjectInfo{
		Name:       "myproject",
		WorkingDir: oldPath,
		ConfigFile: filepath.Join(oldPath, "docker-compose.yml"),
	}

	newDir, newConfig, found := findRelocatedProject(info)
	if !found {
		t.Fatal("移動先が見つからなかった")
	}
	if newDir != newProjectDir {
		t.Errorf("WorkingDir: got %s, want %s", newDir, newProjectDir)
	}
	if newConfig != composeFile {
		t.Errorf("ConfigFile: got %s, want %s", newConfig, composeFile)
	}
}

// パスが変わっていない（存在する）場合はResolveStalePathsがスキップすること
func TestResolveStalePaths_SkipsValidPaths(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "myproject")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}

	registry := &ProjectRegistry{
		projects: map[string]ProjectInfo{
			"myproject": {
				Name:       "myproject",
				WorkingDir: projectDir,
				ConfigFile: filepath.Join(projectDir, "docker-compose.yml"),
			},
		},
	}

	resolved := ResolveStalePaths(registry)
	if len(resolved) != 0 {
		t.Errorf("有効なパスに対して解決が行われた: %v", resolved)
	}
}

// RegisterIfPathExistsが有効な既存パスをstaleパスで上書きしないこと
func TestRegisterIfPathExists_DoesNotOverwriteValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	validDir := filepath.Join(tmpDir, "valid")
	if err := os.MkdirAll(validDir, 0o755); err != nil {
		t.Fatal(err)
	}

	registry := &ProjectRegistry{
		projects: map[string]ProjectInfo{
			"proj": {
				Name:       "proj",
				WorkingDir: validDir,
				ConfigFile: filepath.Join(validDir, "docker-compose.yml"),
			},
		},
	}

	// 存在しないパスで上書きを試みる
	staleDir := filepath.Join(tmpDir, "nonexistent")
	registry.RegisterIfPathExists("proj", staleDir, filepath.Join(staleDir, "docker-compose.yml"))

	info, ok := registry.Get("proj")
	if !ok {
		t.Fatal("プロジェクトが取得できない")
	}
	if info.WorkingDir != validDir {
		t.Errorf("有効なパスが上書きされた: got %s, want %s", info.WorkingDir, validDir)
	}
}

// RegisterIfPathExistsが有効な新パスで登録すること
func TestRegisterIfPathExists_RegistersValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newpath")
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}

	registry := &ProjectRegistry{
		projects: make(map[string]ProjectInfo),
	}

	registry.RegisterIfPathExists("proj", newDir, filepath.Join(newDir, "docker-compose.yml"))

	info, ok := registry.Get("proj")
	if !ok {
		t.Fatal("プロジェクトが登録されていない")
	}
	if info.WorkingDir != newDir {
		t.Errorf("WorkingDir: got %s, want %s", info.WorkingDir, newDir)
	}
}

// 複数候補がある場合、旧パスに近い方が選ばれること
func TestFindRelocatedProject_PrefersClosestMatch(t *testing.T) {
	// tmpdir/
	//   Development/
	//     aaa/
	//       myproject/
	//         docker-compose.yml
	//     bbb/
	//       myproject/
	//         docker-compose.yml
	tmpDir := t.TempDir()
	devDir := filepath.Join(tmpDir, "Development")

	dirA := filepath.Join(devDir, "aaa", "myproject")
	dirB := filepath.Join(devDir, "bbb", "myproject")

	for _, d := range []string{dirA, dirB} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "docker-compose.yml"), []byte("version: '3'"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// 旧パスが Development/aaa/old/myproject だった場合、aaa配下の方が近い
	oldPath := filepath.Join(devDir, "aaa", "old", "myproject")
	info := ProjectInfo{
		Name:       "myproject",
		WorkingDir: oldPath,
		ConfigFile: filepath.Join(oldPath, "docker-compose.yml"),
	}

	newDir, _, found := findRelocatedProject(info)
	if !found {
		t.Fatal("移動先が見つからなかった")
	}
	if newDir != dirA {
		t.Errorf("旧パスに近い候補が選ばれなかった: got %s, want %s", newDir, dirA)
	}
}

// 深さ制限が正しく機能すること
func TestWalkDirLimited_RespectsMaxDepth(t *testing.T) {
	// tmpdir/a/b/c/d/e/target.txt
	tmpDir := t.TempDir()
	deepDir := filepath.Join(tmpDir, "a", "b", "c", "d", "e")
	if err := os.MkdirAll(deepDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deepDir, "target.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 深度2で走査 → a/b/c まで
	var visited []string
	walkDirLimited(tmpDir, 2, func(path string, d os.DirEntry, depth int) error {
		if d.IsDir() {
			visited = append(visited, d.Name())
		}
		return nil
	})

	// d, e は到達しないはず
	for _, name := range visited {
		if name == "d" || name == "e" {
			t.Errorf("深度制限を超えたディレクトリに到達した: %s", name)
		}
	}
}

func TestCommonPrefixLen(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"/Users/foo/Development/proj", "/Users/foo/Development/yy/proj", 4},
		{"/a/b/c", "/a/b/d", 3},  // "" + "a" + "b"
		{"/x/y", "/a/b", 1},      // "" のみ一致
		{"/a", "/a", 2},          // "" + "a"
	}

	for _, tt := range tests {
		got := commonPrefixLen(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("commonPrefixLen(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
