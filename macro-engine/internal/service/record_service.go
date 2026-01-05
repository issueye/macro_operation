package service

import (
	"fmt"
	"sync"

	"macro-engine/internal/engine/events"
	"macro-engine/internal/engine/generator"
	"macro-engine/internal/model"
)

// RecordService 录制服务
type RecordService struct {
	capture   *events.Capture
	generator *generator.Generator
	isStarted bool
	mutex     sync.RWMutex
}

// NewRecordService 创建录制服务
func NewRecordService() *RecordService {
	return &RecordService{
		capture:   events.NewCapture(),
		generator: generator.NewGenerator(),
		isStarted: false,
	}
}

// Start 开始录制
func (s *RecordService) Start(clearPrevious bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isStarted {
		return fmt.Errorf("recording already started")
	}

	if clearPrevious {
		s.capture.ClearEvents()
	}

	if err := s.capture.Start(); err != nil {
		return fmt.Errorf("failed to start capture: %w", err)
	}

	s.isStarted = true
	return nil
}

// Stop 停止录制
func (s *RecordService) Stop() (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isStarted {
		return 0, fmt.Errorf("recording not started")
	}

	s.capture.Stop()
	s.isStarted = false

	return s.capture.GetEventCount(), nil
}

// GetEvents 获取录制的事件
func (s *RecordService) GetEvents() []model.Event {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rawEvents := s.capture.GetEvents()
	events := make([]model.Event, 0, len(rawEvents))

	for _, rawEv := range rawEvents {
		// 转换为 model.Event
		if ev, ok := rawEv.(model.Event); ok {
			events = append(events, ev)
		}
	}

	return events
}

// GenerateScript 生成脚本
func (s *RecordService) GenerateScript(name string) (string, error) {
	events := s.GetEvents()
	if len(events) == 0 {
		return "", fmt.Errorf("no events recorded")
	}

	script, err := s.generator.Generate(name, events)
	if err != nil {
		return "", fmt.Errorf("failed to generate script: %w", err)
	}

	return script, nil
}

// IsStarted 检查是否已开始录制
func (s *RecordService) IsStarted() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isStarted
}

// GetEventCount 获取事件数量
func (s *RecordService) GetEventCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.capture.GetEventCount()
}
