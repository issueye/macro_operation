package events

import (
	"fmt"
	"sync"

	"github.com/robotn/gohook"
)

// Capture 事件捕获器
type Capture struct {
	events   []interface{}
	started  bool
	eventCh  chan interface{}
	stopCh   chan struct{}
	mutex    sync.RWMutex
}

// NewCapture 创建新的捕获器
func NewCapture() *Capture {
	return &Capture{
		events:   make([]interface{}, 0),
		started:  false,
		eventCh:  make(chan interface{}, 1000),
		stopCh:   make(chan struct{}),
	}
}

// Start 开始捕获
func (c *Capture) Start() error {
	c.mutex.Lock()
	if c.started {
		c.mutex.Unlock()
		return fmt.Errorf("capture already started")
	}
	c.started = true
	c.events = make([]interface{}, 0)
	c.stopCh = make(chan struct{})
	c.mutex.Unlock()

	// 启动事件监听
	go c.listenEvents()

	return nil
}

// Stop 停止捕获
func (c *Capture) Stop() {
	c.mutex.Lock()
	if !c.started {
		c.mutex.Unlock()
		return
	}
	c.started = false
	close(c.stopCh)
	c.mutex.Unlock()
}

// Clear 清空事件
func (c *Capture) ClearEvents() {
	c.mutex.Lock()
	c.events = make([]interface{}, 0)
	c.mutex.Unlock()
}

// GetEvents 获取所有事件
func (c *Capture) GetEvents() []interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.events
}

// GetEventCount 获取事件数量
func (c *Capture) GetEventCount() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.events)
}

// listenEvents 监听事件
func (c *Capture) listenEvents() {
	hookEvents := hook.Start()
	defer hook.End()

	for {
		select {
		case <-c.stopCh:
			return
		case ev := <-hookEvents:
			c.mutex.Lock()
			c.events = append(c.events, ev)
			c.mutex.Unlock()
		}
	}
}
