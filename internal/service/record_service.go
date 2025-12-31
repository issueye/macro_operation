package service

import (
	"fmt"
	"sync"

	"github.com/issueye/macro_operation/internal/engine/events"
	"github.com/issueye/macro_operation/internal/engine/generator"
	"github.com/issueye/macro_operation/internal/model"
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
func (s *RecordService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isStarted {
		return fmt.Errorf("recording already started")
	}

	// 清空之前的事件
	s.capture.ClearEvents()

	// 开始捕获
	if err := s.capture.Start(); err != nil {
		return fmt.Errorf("failed to start capture: %w", err)
	}

	s.isStarted = true
	return nil
}

// Stop 停止录制
func (s *RecordService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isStarted {
		return fmt.Errorf("recording not started")
	}

	s.capture.Stop()
	s.isStarted = false

	return nil
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

	// 优化事件
	optimizedEvents := s.generator.OptimizeEvents(events)

	// 生成脚本
	script, err := s.generator.Generate(name, optimizedEvents)
	if err != nil {
		return "", fmt.Errorf("failed to generate script: %w", err)
	}

	return script, nil
}

// GenerateCurrentScript 生成当前已录制事件的脚本
func (s *RecordService) GenerateCurrentScript() (string, error) {
	events := s.GetEvents()
	if len(events) == 0 {
		return "", nil // 没有事件时返回空字符串，不报错
	}

	// 优化事件
	optimizedEvents := s.generator.OptimizeEvents(events)

	// 生成脚本
	// 使用临时名称 "temp"，因为用户还没有保存
	script, err := s.generator.Generate("temp", optimizedEvents)
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
	return s.capture.GetEventCount()
}
