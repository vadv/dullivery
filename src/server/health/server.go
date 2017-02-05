package health

import (
	"golang.org/x/net/context"

	api "api/health"
)

type HealthServer struct {
	ContentDir string
}

func NewHealthServer(contentDir string) *HealthServer {
	return &HealthServer{ContentDir: contentDir}
}

// реализация pong
func (h *HealthServer) Ping(ctx context.Context, ping *api.PingMsg) (*api.PongMsg, error) {
	return &api.PongMsg{}, nil
}
