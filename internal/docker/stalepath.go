package docker

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 標準的なcompose設定ファイル名
var composeFileNames = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

// パス解決の最大探索深度
const maxSearchDepth = 4

// パス解決の結果
type ResolveResult struct {
	ProjectName   string
	OldWorkingDir string
	NewWorkingDir string
	NewConfigFile string
}

// レジストリ内のstaleなパスを検出し、新しいパスへの解決を試みる
func ResolveStalePaths(registry *ProjectRegistry) []ResolveResult {
	allInfos := registry.All()
	var resolved []ResolveResult

	for _, info := range allInfos {
		if info.WorkingDir == "" {
			continue
		}
		// パスが存在する場合はスキップ
		if _, err := os.Stat(info.WorkingDir); err == nil {
			continue
		}

		newDir, newConfig, found := findRelocatedProject(info)
		if !found {
			continue
		}

		registry.Register(info.Name, newDir, newConfig)
		resolved = append(resolved, ResolveResult{
			ProjectName:   info.Name,
			OldWorkingDir: info.WorkingDir,
			NewWorkingDir: newDir,
			NewConfigFile: newConfig,
		})
	}

	return resolved
}

// 移動されたプロジェクトの新しい場所を探索する
func findRelocatedProject(info ProjectInfo) (newWorkingDir, newConfigFile string, found bool) {
	projectDirName := filepath.Base(info.WorkingDir)

	// 探索対象のcomposeファイル名を構築
	targetFiles := buildTargetFileNames(info.ConfigFile)

	// フェーズ1: 旧パスの祖先ディレクトリから探索
	if dir, config, ok := searchFromAncestor(info.WorkingDir, projectDirName, targetFiles); ok {
		return dir, config, true
	}

	// フェーズ2: ホームディレクトリ配下の一般的なパスから探索
	if dir, config, ok := searchFromCommonRoots(projectDirName, targetFiles, info.WorkingDir); ok {
		return dir, config, true
	}

	return "", "", false
}

// 探索対象のcomposeファイル名リストを構築
func buildTargetFileNames(configFile string) []string {
	targets := make([]string, len(composeFileNames))
	copy(targets, composeFileNames)

	if configFile != "" {
		base := filepath.Base(configFile)
		isDuplicate := false
		for _, name := range targets {
			if name == base {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			targets = append(targets, base)
		}
	}

	return targets
}

// 旧パスの祖先ディレクトリを起点に探索
func searchFromAncestor(oldPath, projectDirName string, targetFiles []string) (string, string, bool) {
	// 実在する最も近い祖先ディレクトリを見つける
	anchor := findExistingAncestor(oldPath)
	if anchor == "" || anchor == "/" {
		return "", "", false
	}

	return searchUnderRoot(anchor, projectDirName, targetFiles, oldPath)
}

// ホームディレクトリ配下の一般的なルートから探索
func searchFromCommonRoots(projectDirName string, targetFiles []string, oldPath string) (string, string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", false
	}

	commonRoots := []string{
		filepath.Join(home, "Development"),
		filepath.Join(home, "Projects"),
		filepath.Join(home, "dev"),
		filepath.Join(home, "src"),
		filepath.Join(home, "work"),
	}

	for _, root := range commonRoots {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		if dir, config, ok := searchUnderRoot(root, projectDirName, targetFiles, oldPath); ok {
			return dir, config, true
		}
	}

	return "", "", false
}

// 指定ルート配下で、プロジェクト名に一致するディレクトリ内のcomposeファイルを探索
func searchUnderRoot(root, projectDirName string, targetFiles []string, oldPath string) (string, string, bool) {
	type candidate struct {
		workingDir string
		configFile string
	}
	var candidates []candidate

	walkDirLimited(root, maxSearchDepth, func(path string, d fs.DirEntry, depth int) error {
		if !d.IsDir() {
			return nil
		}
		if d.Name() != projectDirName {
			return nil
		}
		// 旧パスと同じなら除外（存在しないことは確認済み）
		if path == oldPath {
			return nil
		}
		// composeファイルが存在するか確認
		for _, name := range targetFiles {
			configPath := filepath.Join(path, name)
			if _, err := os.Stat(configPath); err == nil {
				candidates = append(candidates, candidate{
					workingDir: path,
					configFile: configPath,
				})
				return fs.SkipDir
			}
		}
		return nil
	})

	if len(candidates) == 0 {
		return "", "", false
	}

	// 旧パスとの共通プレフィックスが最長の候補を選択
	best := candidates[0]
	bestPrefix := commonPrefixLen(oldPath, best.workingDir)
	for _, c := range candidates[1:] {
		prefix := commonPrefixLen(oldPath, c.workingDir)
		if prefix > bestPrefix {
			best = c
			bestPrefix = prefix
		}
	}

	return best.workingDir, best.configFile, true
}

// 実在する最も近い祖先ディレクトリを返す
func findExistingAncestor(path string) string {
	dir := filepath.Dir(path)
	for dir != "/" && dir != "." {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return dir
}

// 深さ制限付きディレクトリ走査
func walkDirLimited(root string, maxDepth int, fn func(path string, d fs.DirEntry, depth int) error) {
	walkRecursive(root, 0, maxDepth, fn)
}

func walkRecursive(dir string, currentDepth, maxDepth int, fn func(path string, d fs.DirEntry, depth int) error) {
	if currentDepth > maxDepth {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		err := fn(path, entry, currentDepth)
		if err == fs.SkipDir {
			continue
		}
		if err == fs.SkipAll {
			return
		}
		if entry.IsDir() && currentDepth < maxDepth {
			// 隠しディレクトリとnode_modules等はスキップ
			name := entry.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				continue
			}
			walkRecursive(path, currentDepth+1, maxDepth, fn)
		}
	}
}

// 2つのパスの共通プレフィックスの長さを返す
func commonPrefixLen(a, b string) int {
	partsA := strings.Split(a, string(filepath.Separator))
	partsB := strings.Split(b, string(filepath.Separator))

	n := len(partsA)
	if len(partsB) < n {
		n = len(partsB)
	}

	common := 0
	for i := 0; i < n; i++ {
		if partsA[i] != partsB[i] {
			break
		}
		common++
	}
	return common
}
