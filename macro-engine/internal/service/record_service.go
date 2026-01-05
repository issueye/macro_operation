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

	// 获取所有事件用于诊断
	events := s.capture.GetEvents()

	// 打印诊断信息
	fmt.Printf("[DIAGNOSTIC] Total events captured: %d\n", len(events))
	keyDownCount := 0
	keyUpCount := 0
	for i, ev := range events {
		if ev.Type == model.EventTypeKeyDown {
			keyDownCount++
		} else if ev.Type == model.EventTypeKeyUp {
			keyUpCount++
		}
		// 打印前 50 个事件的详情
		if i < 50 {
			fmt.Printf("[DIAGNOSTIC] Event %d: Type=%s, KeyCode=%d, Chars=%q, X=%d, Y=%d\n",
				i, ev.Type, ev.KeyCode, ev.Chars, ev.X, ev.Y)
		}
	}
	fmt.Printf("[DIAGNOSTIC] KeyDown events: %d, KeyUp events: %d\n", keyDownCount, keyUpCount)

	return len(events), nil
}

// GetEvents 获取录制的事件
func (s *RecordService) GetEvents() []model.Event {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.capture.GetEvents()
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
