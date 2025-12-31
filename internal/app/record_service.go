package app

import (
	"context"
	"sync"
	"time"

	"github.com/issueye/macro_operation/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// RecordService Wails v3 录制服务
type RecordService struct {
	service   *service.RecordService
	app       *application.App
	mutex     sync.RWMutex
	startTime time.Time
}

// RecordingStatus 录制状态
type RecordingStatus struct {
	IsStarted  bool   `json:"is_started"`
	EventCount int    `json:"event_count"`
	DurationMs int64  `json:"duration_ms"`
	StartTime  string `json:"start_time"`
}

// NewRecordService 创建录制服务
func NewRecordService(app *application.App) *RecordService {
	return &RecordService{
		service: service.NewRecordService(),
		app:     app,
	}
}

// Start 开始录制
func (s *RecordService) Start(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.service.Start(); err != nil {
		return err
	}

	s.startTime = time.Now()

	// TODO: 发送事件 - 需要查看 Wails v3 正确的事件 API
	// events.Emit("recording:started", nil)

	// 启动状态更新
	go s.sendStatusUpdates()

	return nil
}

// Stop 停止录制
func (s *RecordService) Stop(ctx context.Context) (*RecordingStatus, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.service.Stop(); err != nil {
		return nil, err
	}

	status := s.getStatus()
	// TODO: 发送事件
	// events.Emit("recording:stopped", status)

	return status, nil
}

// GetStatus 获取录制状态
func (s *RecordService) GetStatus(ctx context.Context) (*RecordingStatus, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.getStatus(), nil
}

// GenerateScript 生成脚本
func (s *RecordService) GenerateScript(ctx context.Context, name string) (string, error) {
	return s.service.GenerateScript(name)
}

// 内部方法
func (s *RecordService) getStatus() *RecordingStatus {
	duration := time.Duration(0)
	if !s.startTime.IsZero() && s.service.IsStarted() {
		duration = time.Since(s.startTime)
	}

	return &RecordingStatus{
		IsStarted:  s.service.IsStarted(),
		EventCount: s.service.GetEventCount(),
		DurationMs: duration.Milliseconds(),
		StartTime:  s.startTime.Format("2006-01-02 15:04:05"),
	}
}

func (s *RecordService) sendStatusUpdates() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for s.service.IsStarted() {
		select {
		case <-ticker.C:
			status := s.getStatus()
			// TODO: 发送进度事件
			_ = status // 暂时不发送事件
			// events.Emit("recording:progress", status)
		}
	}
}
