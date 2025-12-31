package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
	"github.com/issueye/macro_operation/internal/model"
)

// Capture 事件捕获器
type Capture struct {
	keyboardChan chan model.KeyboardEvent
	mouseChan    chan model.MouseEvent
	isRunning    bool
	eventBuffer  []model.Event
	bufferMutex  sync.RWMutex
	wg           sync.WaitGroup
	stopChan     chan struct{}
}

// NewCapture 创建事件捕获器
func NewCapture() *Capture {
	return &Capture{
		keyboardChan: make(chan model.KeyboardEvent, 1000),
		mouseChan:    make(chan model.MouseEvent, 1000),
		isRunning:    false,
		eventBuffer:  make([]model.Event, 0, 10000),
		stopChan:     make(chan struct{}),
	}
}

// Start 开始捕获事件
func (c *Capture) Start() error {
	if c.isRunning {
		return fmt.Errorf("capture already started")
	}

	c.isRunning = true
	c.eventBuffer = make([]model.Event, 0, 10000)

	// 启动键盘监听
	c.wg.Add(1)
	go c.startKeyboardHook()

	// 启动鼠标监听
	c.wg.Add(1)
	go c.startMouseHook()

	// 启动事件处理
	c.wg.Add(1)
	go c.processEvents()

	return nil
}

// startKeyboardHook 启动键盘钩子
func (c *Capture) startKeyboardHook() {
	defer c.wg.Done()

	evChan := hook.Start()

	for {
		select {
		case <-c.stopChan:
			return
		case ev := <-evChan:
			if !c.isRunning {
				return
			}

			switch ev.Kind {
			case hook.KeyDown:
				c.keyboardChan <- model.KeyboardEvent{
					Type:      model.KeyDown,
					Keycode:   ev.Keychar,
					Keyname:   string(rune(ev.Keychar)),
					Timestamp: time.Now(),
				}
			case hook.KeyUp:
				c.keyboardChan <- model.KeyboardEvent{
					Type:      model.KeyUp,
					Keycode:   ev.Keychar,
					Keyname:   string(rune(ev.Keychar)),
					Timestamp: time.Now(),
				}
			case hook.KeyHold:
				c.keyboardChan <- model.KeyboardEvent{
					Type:      model.KeyPress,
					Keycode:   ev.Keychar,
					Keyname:   string(rune(ev.Keychar)),
					Timestamp: time.Now(),
				}
			}
		}
	}
}

// startMouseHook 启动鼠标监听
func (c *Capture) startMouseHook() {
	defer c.wg.Done()

	var lastX, lastY int

	for c.isRunning {
		select {
		case <-c.stopChan:
			return
		default:
			x, y := robotgo.GetMousePos()

			// 检查鼠标移动
			if x != lastX || y != lastY {
				c.mouseChan <- model.MouseEvent{
					Type:      model.MouseMove,
					X:         x,
					Y:         y,
					Timestamp: time.Now(),
				}
				lastX, lastY = x, y
			}

			time.Sleep(20 * time.Millisecond)
		}
	}
}

// processEvents 处理事件流
func (c *Capture) processEvents() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		case ke := <-c.keyboardChan:
			c.bufferMutex.Lock()
			c.eventBuffer = append(c.eventBuffer, ke)
			c.bufferMutex.Unlock()

		case me := <-c.mouseChan:
			c.bufferMutex.Lock()
			c.eventBuffer = append(c.eventBuffer, me)
			c.bufferMutex.Unlock()
		}
	}
}

// Stop 停止捕获
func (c *Capture) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false
	close(c.stopChan)
	hook.End()
	c.wg.Wait()

	// 重新创建 stopChan
	c.stopChan = make(chan struct{})
}

// GetEvents 获取捕获的事件
func (c *Capture) GetEvents() []model.Event {
	c.bufferMutex.RLock()
	defer c.bufferMutex.RUnlock()

	events := make([]model.Event, len(c.eventBuffer))
	copy(events, c.eventBuffer)
	return events
}

// ClearEvents 清空事件缓冲
func (c *Capture) ClearEvents() {
	c.bufferMutex.Lock()
	defer c.bufferMutex.Unlock()

	c.eventBuffer = make([]model.Event, 0, 10000)
}

// IsRunning 检查是否正在运行
func (c *Capture) IsRunning() bool {
	return c.isRunning
}

// GetEventCount 获取事件数量
func (c *Capture) GetEventCount() int {
	c.bufferMutex.RLock()
	defer c.bufferMutex.RUnlock()
	return len(c.eventBuffer)
}

// FilterEvents 过滤事件
func (c *Capture) FilterEvents(filterFn func(model.Event) bool) []model.Event {
	c.bufferMutex.RLock()
	defer c.bufferMutex.RUnlock()

	var filtered []model.Event
	for _, event := range c.eventBuffer {
		if filterFn(event) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}
