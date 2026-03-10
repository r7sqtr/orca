package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/r7sqtr/orca/internal/i18n"
	"github.com/r7sqtr/orca/internal/model"
	"github.com/r7sqtr/orca/internal/ui"
)

// サービス詳細パネル
type ServiceInfo struct {
	styles    ui.Styles
	width     int
	height    int
	container *model.ContainerStatus
	project   string
	service   string
}

// ServiceInfoを作成
func NewServiceInfo(styles ui.Styles) ServiceInfo {
	return ServiceInfo{styles: styles}
}

// サイズを設定
func (si *ServiceInfo) SetSize(width, height int) {
	si.width = width
	si.height = height
}

// 表示するサービスを設定
func (si *ServiceInfo) SetService(project, service string, container *model.ContainerStatus) {
	si.project = project
	si.service = service
	si.container = container
}

// サービス情報を描画
func (si ServiceInfo) View() string {
	if si.container == nil {
		if si.service != "" {
			return si.styles.Muted.Render(
				fmt.Sprintf("%s: %s (%s)", i18n.T("detail.service"), si.service, i18n.T("status.not_created")),
			)
		}
		return ""
	}

	c := si.container
	var lines []string

	// ラベルと値のペア
	pairs := []struct {
		label string
		value string
		style lipgloss.Style
	}{
		{i18n.T("detail.service"), si.service, si.styles.Bold},
		{i18n.T("detail.state"), si.stateText(c.State), si.stateStyle(c.State)},
		{i18n.T("detail.id"), c.ID, si.styles.Muted},
		{i18n.T("detail.image"), c.Image, lipgloss.NewStyle()},
	}

	// ポート
	if len(c.Ports) > 0 {
		portStrs := make([]string, 0, len(c.Ports))
		for _, p := range c.Ports {
			if p.HostPort > 0 {
				portStrs = append(portStrs, fmt.Sprintf("%s:%d->%d/%s",
					p.HostIP, p.HostPort, p.ContainerPort, p.Protocol))
			} else {
				portStrs = append(portStrs, fmt.Sprintf("%d/%s",
					p.ContainerPort, p.Protocol))
			}
		}
		pairs = append(pairs, struct {
			label string
			value string
			style lipgloss.Style
		}{i18n.T("detail.ports"), strings.Join(portStrs, ", "), lipgloss.NewStyle()})
	}

	// ヘルスチェック
	if c.Health != "" {
		pairs = append(pairs, struct {
			label string
			value string
			style lipgloss.Style
		}{i18n.T("detail.health"), c.Health, si.healthStyle(c.Health)})
	}

	labelWidth := 0
	for _, p := range pairs {
		w := lipgloss.Width(p.label)
		if w > labelWidth {
			labelWidth = w
		}
	}

	for _, p := range pairs {
		label := si.styles.Subtitle.Width(labelWidth + 1).Render(p.label + ":")
		value := p.style.Render(p.value)
		lines = append(lines, label+" "+value)
	}

	return strings.Join(lines, "\n")
}

func (si ServiceInfo) stateText(state string) string {
	key := "status." + state
	text := i18n.T(key)
	if text == key {
		return state
	}
	return text
}

func (si ServiceInfo) stateStyle(state string) lipgloss.Style {
	switch state {
	case "running":
		return si.styles.Running
	case "paused":
		return si.styles.Warning
	case "exited", "dead":
		return si.styles.Stopped
	case "restarting":
		return si.styles.Warning
	default:
		return si.styles.Muted
	}
}

func (si ServiceInfo) healthStyle(health string) lipgloss.Style {
	switch health {
	case "healthy":
		return si.styles.Health
	case "unhealthy":
		return si.styles.Error
	default:
		return si.styles.Warning
	}
}
