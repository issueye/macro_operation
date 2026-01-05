package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

// ClipboardMonitor 剪贴板监控器 - 用于捕获中文输入
// 工作原理：当用户通过输入法输入中文时，通常会经过剪贴板或可以直接获取
type ClipboardMonitor struct {
	lastContent string
	eventCh     chan string
	stopCh      chan struct{}
	started     bool
	mutex       sync.RWMutex
}

// NewClipboardMonitor 创建剪贴板监控器
func NewClipboardMonitor() *ClipboardMonitor {
	return &ClipboardMonitor{
		eventCh: make(chan string, 100),
		stopCh:  make(chan struct{}),
		started: false,
	}
}

// Start 开始监控剪贴板
func (m *ClipboardMonitor) Start() error {
	m.mutex.Lock()
	if m.started {
		m.mutex.Unlock()
		return fmt.Errorf("clipboard monitor already started")
	}
	m.started = true
	m.stopCh = make(chan struct{})

	// 获取初始剪贴板内容
	initialContent, err := clipboard.ReadAll()
	if err == nil {
		m.lastContent = initialContent
	}

	m.mutex.Unlock()

	// 启动监控协程
	go m.monitor()

	return nil
}

// Stop 停止监控
func (m *ClipboardMonitor) Stop() {
	m.mutex.Lock()
	if !m.started {
		m.mutex.Unlock()
		return
	}
	m.started = false
	close(m.stopCh)
	m.mutex.Unlock()
}

// monitor 监控剪贴板变化
func (m *ClipboardMonitor) monitor() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			current, err := clipboard.ReadAll()
			if err != nil {
				continue
			}

			m.mutex.Lock()
			if current != m.lastContent && current != "" {
				// 剪贴板内容发生了变化
				// 检查是否包含非 ASCII 字符（可能是中文输入）
				if hasNonASCII(current) {
					// 发送字符事件
					select {
					case m.eventCh <- current:
					default:
						// 通道满了，丢弃
					}
				}
				m.lastContent = current
			}
			m.mutex.Unlock()
		}
	}
}

// GetEventChannel 获取事件通道
func (m *ClipboardMonitor) GetEventChannel() <-chan string {
	return m.eventCh
}

// hasNonASCII 检查字符串是否包含非 ASCII 字符
func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}
