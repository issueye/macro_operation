package service

import (
	"fmt"
	"sync"

	"macro-engine/internal/engine/executor"
)

// PlayService 回放服务
type PlayService struct {
	executor  *executor.Executor
	isPlaying bool
	mutex     sync.RWMutex
}

// NewPlayService 创建回放服务
func NewPlayService() *PlayService {
	return &PlayService{
		executor:  executor.NewExecutor(),
		isPlaying: false,
	}
}

// PlayScript 执行脚本
func (s *PlayService) PlayScript(script string, speed float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isPlaying {
		return fmt.Errorf("playback already in progress")
	}

	s.isPlaying = true
	defer func() {
		s.isPlaying = false
	}()

	// TODO: 根据 speed 调整执行速度
	_ = speed

	if err := s.executor.Execute(script); err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return nil
}

// IsPlaying 检查是否正在播放
func (s *PlayService) IsPlaying() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isPlaying
}
