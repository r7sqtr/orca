package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	dockerclient "github.com/docker/docker/client"
	"github.com/r7sqtr/orca/internal/model"
)

// Docker APIクライアント
type Client struct {
	cli *dockerclient.Client
}

// Colima環境のDockerソケットを自動検出
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

// Dockerクライアントを作成
func NewClient(dockerHost string) (*Client, error) {
	// 設定の DockerHost を最優先で使用
	if dockerHost != "" {
		os.Setenv("DOCKER_HOST", dockerHost)
	} else if host := detectDockerHost(); host != "" {
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

// Dockerデーモンへの接続を確認
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

// クライアントを閉じる
func (c *Client) Close() error {
	return c.cli.Close()
}

// 全コンテナを取得
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

// Composeプロジェクトのコンテナのみ取得
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

// Dockerイベントストリームを返却
func (c *Client) Events(ctx context.Context) (<-chan events.Message, <-chan error) {
	return c.cli.Events(ctx, events.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
		),
	})
}

// コンテナの詳細情報を取得
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

// コンテナの環境変数を取得
func (c *Client) GetContainerEnv(ctx context.Context, containerID string) ([]string, error) {
	info, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return info.Config.Env, nil
}

// コンテナとcompose設定のサービスをマージしてグループ化
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

// 全イメージを取得
func (c *Client) ListImages(ctx context.Context) ([]model.ImageInfo, error) {
	images, err := c.cli.ImageList(ctx, image.ListOptions{
		All:            true,
		ContainerCount: true,
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.ImageInfo, 0, len(images))
	for _, img := range images {
		info := model.ImageInfo{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
		}
		if img.Created > 0 {
			info.Created = time.Unix(img.Created, 0)
		}
		if img.Containers >= 0 {
			info.Containers = img.Containers
		}
		result = append(result, info)
	}

	return result, nil
}

// イメージを削除
func (c *Client) RemoveImage(ctx context.Context, imageID string) error {
	_, err := c.cli.ImageRemove(ctx, imageID, image.RemoveOptions{
		PruneChildren: true,
	})
	return err
}

// 未使用イメージを一括削除
func (c *Client) PruneImages(ctx context.Context) (uint64, error) {
	report, err := c.cli.ImagesPrune(ctx, filters.NewArgs(
		filters.Arg("dangling", "false"),
	))
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// 全ボリュームを取得
func (c *Client) ListVolumes(ctx context.Context) ([]model.VolumeInfo, error) {
	resp, err := c.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]model.VolumeInfo, 0, len(resp.Volumes))
	for _, vol := range resp.Volumes {
		info := model.VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			MountPoint: vol.Mountpoint,
			Labels:     vol.Labels,
		}
		if vol.CreatedAt != "" {
			if t, err := time.Parse(time.RFC3339Nano, vol.CreatedAt); err == nil {
				info.CreatedAt = t
			}
		}
		if vol.UsageData != nil {
			info.RefCount = vol.UsageData.RefCount
		}
		result = append(result, info)
	}

	return result, nil
}

// ボリュームを削除
func (c *Client) RemoveVolume(ctx context.Context, volumeName string) error {
	return c.cli.VolumeRemove(ctx, volumeName, false)
}

// 未使用ボリュームを一括削除
func (c *Client) PruneVolumes(ctx context.Context) (uint64, error) {
	report, err := c.cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// ステータス文字列からヘルス状態を抽出
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
