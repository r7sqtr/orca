package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ProjectInfo はプロジェクトの設定情報を保持する
type ProjectInfo struct {
	Name       string `json:"name"`
	WorkingDir string `json:"working_dir"`
	ConfigFile string `json:"config_file"`
}

// ProjectRegistry はプロジェクト情報を永続的に管理する
type ProjectRegistry struct {
	mu       sync.RWMutex
	projects map[string]ProjectInfo
	filePath string
	dirty    bool
}

// NewProjectRegistry はProjectRegistryを作成し、保存済みデータを読み込む
func NewProjectRegistry() *ProjectRegistry {
	r := &ProjectRegistry{
		projects: make(map[string]ProjectInfo),
		filePath: registryPath(),
	}
	r.load()
	return r
}

// Register はプロジェクト情報を登録する（保存はSave()で明示的に行う）
func (r *ProjectRegistry) Register(name, workingDir, configFile string) {
	if name == "" || workingDir == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.projects[name]
	if ok && existing.WorkingDir == workingDir && existing.ConfigFile == configFile {
		return
	}

	r.projects[name] = ProjectInfo{
		Name:       name,
		WorkingDir: workingDir,
		ConfigFile: configFile,
	}
	r.dirty = true
}

// All は登録済み全プロジェクト情報を返す
func (r *ProjectRegistry) All() []ProjectInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ProjectInfo, 0, len(r.projects))
	for _, p := range r.projects {
		result = append(result, p)
	}
	return result
}

// Remove はプロジェクトをレジストリから削除する
func (r *ProjectRegistry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.projects[name]; ok {
		delete(r.projects, name)
		r.dirty = true
	}
}

// Save は変更がある場合にファイルに永続化する
func (r *ProjectRegistry) Save() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.dirty {
		return nil
	}
	return r.saveLocked()
}

func registryPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "orca", "registry.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "orca", "registry.json")
}

func (r *ProjectRegistry) load() {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}
	var projects []ProjectInfo
	if err := json.Unmarshal(data, &projects); err != nil {
		return
	}
	for _, p := range projects {
		r.projects[p.Name] = p
	}
}

// saveLocked はロック取得済みの状態でファイルに保存する
func (r *ProjectRegistry) saveLocked() error {
	projects := make([]ProjectInfo, 0, len(r.projects))
	for _, p := range r.projects {
		projects = append(projects, p)
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("レジストリのJSON変換に失敗: %w", err)
	}
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("レジストリディレクトリの作成に失敗: %w", err)
	}
	if err := os.WriteFile(r.filePath, data, 0o644); err != nil {
		return fmt.Errorf("レジストリの保存に失敗: %w", err)
	}
	r.dirty = false
	return nil
}
