package service

import (
	"fmt"
	"sync"

	"macro-engine/internal/engine/events"
	"macro-engine/internal/engine/generator"
	"macro-engine/internal/model"
)

// HotkeyAction 快捷键动作类型
type HotkeyAction string

const (
	HotkeyActionToggleRecording HotkeyAction = "toggle_recording" // 切换录制状态
	HotkeyActionSave            HotkeyAction = "save"             // 保存
	HotkeyActionClear           HotkeyAction = "clear"            // 清空
)

// RecordService 录制服务
type RecordService struct {
	capture        *events.Capture
	generator      *generator.Generator
	isStarted      bool
	mutex          sync.RWMutex
	actionCallback map[HotkeyAction]func() // 快捷键动作回调
}

// NewRecordService 创建录制服务
func NewRecordService() *RecordService {
	rs := &RecordService{
		capture:        events.NewCapture(),
		generator:      generator.NewGenerator(),
		isStarted:      false,
		actionCallback: make(map[HotkeyAction]func()),
	}

	// TODO: 重新启用快捷键功能
	// rs.registerDefaultHotkeys()

	return rs
}

// registerDefaultHotkeys 注册默认快捷键
func (s *RecordService) registerDefaultHotkeys() {
	// 快捷键功能暂时禁用，等待 capture.go 更新
	/*
	// Ctrl+R: 切换录制状态
	s.capture.RegisterHotkey("ctrl+r", func() {
		s.handleToggleRecording()
	})

	// Ctrl+S: 保存
	s.capture.RegisterHotkey("ctrl+s", func() {
		s.handleSave()
	})

	// Ctrl+Shift+C: 清空录制
	s.capture.RegisterHotkey("ctrl+shift+c", func() {
		s.handleClear()
	})

	// Ctrl+Shift+R: 停止录制并保存
	s.capture.RegisterHotkey("ctrl+shift+r", func() {
		s.handleStopAndSave()
	})
	*/
}

// handleToggleRecording 处理切换录制状态
func (s *RecordService) handleToggleRecording() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isStarted {
		fmt.Println("[快捷键] 停止录制")
		s.capture.Stop()
		s.isStarted = false
		eventCount := len(s.capture.GetEvents())
		fmt.Printf("[快捷键] 录制已停止，共捕获 %d 个事件\n", eventCount)
	} else {
		fmt.Println("[快捷键] 开始录制")
		s.capture.ClearEvents()
		if err := s.capture.Start(); err != nil {
			fmt.Printf("[快捷键] 启动录制失败: %v\n", err)
			return
		}
		s.isStarted = true
		fmt.Println("[快捷键] 录制已开始")
	}

	// 触发回调
	if callback, ok := s.actionCallback[HotkeyActionToggleRecording]; ok {
		go callback()
	}
}

// handleSave 处理保存
func (s *RecordService) handleSave() {
	s.mutex.RLock()
	eventCount := len(s.capture.GetEvents())
	s.mutex.RUnlock()

	fmt.Printf("[快捷键] 保存 - 当前有 %d 个事件\n", eventCount)

	// 触发回调
	if callback, ok := s.actionCallback[HotkeyActionSave]; ok {
		go callback()
	}
}

// handleClear 处理清空
func (s *RecordService) handleClear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Println("[快捷键] 清空录制")
	s.capture.ClearEvents()
	fmt.Println("[快捷键] 事件已清空")

	// 触发回调
	if callback, ok := s.actionCallback[HotkeyActionClear]; ok {
		go callback()
	}
}

// handleStopAndSave 处理停止录制并保存
func (s *RecordService) handleStopAndSave() {
	s.mutex.Lock()
	if s.isStarted {
		s.capture.Stop()
		s.isStarted = false
	}
	s.mutex.Unlock()

	eventCount := len(s.capture.GetEvents())
	fmt.Printf("[快捷键] 停止录制并保存，共 %d 个事件\n", eventCount)

	// 触发保存回调
	if callback, ok := s.actionCallback[HotkeyActionSave]; ok {
		go callback()
	}
}

// SetActionCallback 设置快捷键动作的回调函数
func (s *RecordService) SetActionCallback(action HotkeyAction, callback func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.actionCallback[action] = callback
}

// GetCapture 获取捕获器实例（用于外部注册自定义快捷键）
func (s *RecordService) GetCapture() *events.Capture {
	return s.capture
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

	// 启动事件捕获（如果还没启动）
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
