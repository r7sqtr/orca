package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/vvsaito/orca/internal/model"
)

// Client はDocker APIクライアント
type Client struct {
	cli *dockerclient.Client
}

// detectDockerHost はColima環境のDockerソケットを自動検出する
func detectDockerHost() string {
	// 既にDOCKER_HOSTが設定されている場合はそのまま使用
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		return host
	}

	// デフォルトソケットが存在する場合
	if _, err := os.Stat("/var/run/docker.sock"); err == nil {
		return ""
	}

	// Colimaのソケットを探索
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	colimaSocket := filepath.Join(home, ".colima", "default", "docker.sock")
	if _, err := os.Stat(colimaSocket); err == nil {
		return "unix://" + colimaSocket
	}

	return ""
}

// NewClient はDockerクライアントを作成する
func NewClient() (*Client, error) {
	// Colima環境のソケットを自動検出
	if host := detectDockerHost(); host != "" {
		os.Setenv("DOCKER_HOST", host)
	}

	cli, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("Docker クライアントの作成に失敗: %w", err)
	}
	return &Client{cli: cli}, nil
}

// Ping はDockerデーモンへの接続を確認する
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

// Close はクライアントを閉じる
func (c *Client) Close() error {
	return c.cli.Close()
}

// ListContainers は全コンテナを取得する
func (c *Client) ListContainers(ctx context.Context) ([]model.ContainerStatus, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ContainerStatus, 0, len(containers))
	for _, ctr := range containers {
		cs := model.ContainerStatus{
			ID:     ctr.ID[:12],
			Name:   ctr.Names[0],
			Image:  ctr.Image,
			State:  ctr.State,
			Status: ctr.Status,
			Labels: ctr.Labels,
		}

		// ポートマッピング
		for _, p := range ctr.Ports {
			cs.Ports = append(cs.Ports, model.PortMapping{
				HostIP:        p.IP,
				HostPort:      p.PublicPort,
				ContainerPort: p.PrivatePort,
				Protocol:      p.Type,
			})
		}

		// Composeラベル
		if proj, ok := ctr.Labels[LabelComposeProject]; ok {
			cs.ProjectName = proj
		}
		if svc, ok := ctr.Labels[LabelComposeService]; ok {
			cs.ServiceName = svc
		}

		result = append(result, cs)
	}

	return result, nil
}

// ListComposeContainers はComposeプロジェクトのコンテナのみ取得する
func (c *Client) ListComposeContainers(ctx context.Context) ([]model.ContainerStatus, error) {
	f := filters.NewArgs()
	f.Add("label", LabelComposeProject)

	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ContainerStatus, 0, len(containers))
	for _, ctr := range containers {
		name := ""
		if len(ctr.Names) > 0 {
			name = ctr.Names[0]
		}

		cs := model.ContainerStatus{
			ID:          ctr.ID[:12],
			Name:        name,
			Image:       ctr.Image,
			State:       ctr.State,
			Status:      ctr.Status,
			Labels:      ctr.Labels,
			ProjectName: ctr.Labels[LabelComposeProject],
			ServiceName: ctr.Labels[LabelComposeService],
		}

		for _, p := range ctr.Ports {
			cs.Ports = append(cs.Ports, model.PortMapping{
				HostIP:        p.IP,
				HostPort:      p.PublicPort,
				ContainerPort: p.PrivatePort,
				Protocol:      p.Type,
			})
		}

		// ヘルスチェック
		if ctr.State == "running" {
			cs.Health = extractHealth(ctr.Status)
		}

		result = append(result, cs)
	}

	return result, nil
}

// Events はDockerイベントストリームを返す
func (c *Client) Events(ctx context.Context) (<-chan events.Message, <-chan error) {
	return c.cli.Events(ctx, events.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
		),
	})
}

// InspectContainer はコンテナの詳細情報を取得する
func (c *Client) InspectContainer(ctx context.Context, containerID string) (*model.ContainerStatus, error) {
	info, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	cs := &model.ContainerStatus{
		ID:     info.ID[:12],
		Name:   info.Name,
		Image:  info.Config.Image,
		State:  info.State.Status,
		Status: info.State.Status,
		Labels: info.Config.Labels,
	}

	if proj, ok := info.Config.Labels[LabelComposeProject]; ok {
		cs.ProjectName = proj
	}
	if svc, ok := info.Config.Labels[LabelComposeService]; ok {
		cs.ServiceName = svc
	}

	if info.State.Health != nil {
		cs.Health = string(info.State.Health.Status)
	}

	return cs, nil
}

// GetContainerEnv はコンテナの環境変数を取得する
func (c *Client) GetContainerEnv(ctx context.Context, containerID string) ([]string, error) {
	info, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return info.Config.Env, nil
}

// GroupByProjectWithConfig はコンテナとcompose設定のサービスをマージしてグループ化する
func GroupByProjectWithConfig(containers []model.ContainerStatus, configServices map[string][]string) []model.ComposeProject {
	projectMap := make(map[string]*model.ComposeProject)
	serviceMap := make(map[string]map[string]*model.Service)

	// コンテナからプロジェクト・サービスを構築
	for _, ctr := range containers {
		projName := ctr.ProjectName
		if projName == "" {
			continue
		}

		if _, ok := projectMap[projName]; !ok {
			projectMap[projName] = &model.ComposeProject{
				Name:       projName,
				WorkingDir: ctr.Labels[LabelComposeWorkingDir],
				ConfigFile: ctr.Labels[LabelComposeConfigFile],
			}
			serviceMap[projName] = make(map[string]*model.Service)
		}

		svcName := ctr.ServiceName
		if _, ok := serviceMap[projName][svcName]; !ok {
			c := ctr
			serviceMap[projName][svcName] = &model.Service{
				Name:        svcName,
				ProjectName: projName,
				Container:   &c,
			}
		}
	}

	// compose設定のサービスをマージ（コンテナが存在しないサービスを追加）
	for projName, svcNames := range configServices {
		if _, ok := projectMap[projName]; !ok {
			projectMap[projName] = &model.ComposeProject{Name: projName}
			serviceMap[projName] = make(map[string]*model.Service)
		}
		for _, svc := range svcNames {
			if _, ok := serviceMap[projName][svc]; !ok {
				serviceMap[projName][svc] = &model.Service{
					Name:        svc,
					ProjectName: projName,
					Container:   nil,
				}
			}
		}
	}

	// スライスに変換してソート
	projects := make([]model.ComposeProject, 0, len(projectMap))
	for name, proj := range projectMap {
		services := make([]model.Service, 0, len(serviceMap[name]))
		for _, svc := range serviceMap[name] {
			services = append(services, *svc)
		}
		sort.Slice(services, func(i, j int) bool {
			return services[i].Name < services[j].Name
		})
		proj.Services = services
		projects = append(projects, *proj)
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	return projects
}

// extractHealth はステータス文字列からヘルス状態を抽出する
func extractHealth(status string) string {
	// "Up 2 hours (healthy)" -> "healthy"
	for _, h := range []string{"healthy", "unhealthy", "starting"} {
		if contains(status, h) {
			return h
		}
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
