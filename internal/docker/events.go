package docker

import (
	"context"
)

// Dockerイベントの種別
type EventType string

const (
	EventStart   EventType = "start"
	EventStop    EventType = "stop"
	EventDie     EventType = "die"
	EventCreate  EventType = "create"
	EventDestroy EventType = "destroy"
	EventRestart EventType = "restart"
	EventPause   EventType = "pause"
	EventUnpause EventType = "unpause"
	EventHealth  EventType = "health_status"
)

// UIに送信されるDockerイベント
type DockerEvent struct {
	Type        EventType
	ContainerID string
	Service     string
	Project     string
}

// Dockerイベントを監視し、チャネルに送信
func WatchEvents(ctx context.Context, client *Client, ch chan<- DockerEvent) {
	msgCh, errCh := client.Events(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgCh:
			event := DockerEvent{
				Type:        EventType(msg.Action),
				ContainerID: msg.Actor.ID[:12],
			}
			if proj, ok := msg.Actor.Attributes[LabelComposeProject]; ok {
				event.Project = proj
			}
			if svc, ok := msg.Actor.Attributes[LabelComposeService]; ok {
				event.Service = svc
			}

			select {
			case ch <- event:
			default:
			}
		case <-errCh:
			return
		}
	}
}
