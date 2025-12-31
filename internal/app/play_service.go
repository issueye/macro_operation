package app

import (
	"context"

	"github.com/issueye/macro_operation/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// PlayService Wails v3 播放服务
type PlayService struct {
	service *service.PlayService
	app     *application.App
}

// PlaybackStatus 播放状态
type PlaybackStatus struct {
	IsPlaying bool   `json:"is_playing"`
	Progress  int    `json:"progress"` // 0-100
	Error     string `json:"error,omitempty"`
}

// NewPlayService 创建播放服务
func NewPlayService(app *application.App) *PlayService {
	return &PlayService{
		service: service.NewPlayService(),
		app:     app,
	}
}

// PlayScript 播放脚本
func (s *PlayService) PlayScript(ctx context.Context, script string) error {
	// TODO: 发送事件
	// events.Emit("playback:started", nil)

	go func() {
		if err := s.service.PlayScript(script); err != nil {
			// TODO: 发送错误事件
			_ = err
			// events.Emit("playback:error", map[string]string{
			//     "error": err.Error(),
			// })
			return
		}
		// TODO: 发送完成事件
		// events.Emit("playback:completed", nil)
	}()

	return nil
}

// IsPlaying 检查是否正在播放
func (s *PlayService) IsPlaying(ctx context.Context) (bool, error) {
	return s.service.IsPlaying(), nil
}

// Stop 停止播放
func (s *PlayService) Stop(ctx context.Context) error {
	// TODO: 发送事件
	// events.Emit("playback:stopped", nil)
	return nil
}
