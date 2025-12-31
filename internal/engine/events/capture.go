package events

import (
	"fmt"
	"sync"
	"time"

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

	// 启动事件监听（统一处理键盘和鼠标事件）
	c.wg.Add(1)
	go c.startEventCapture()

	// 启动事件处理
	c.wg.Add(1)
	go c.processEvents()

	return nil
}

// startEventCapture 统一处理键盘和鼠标事件
func (c *Capture) startEventCapture() {
	defer c.wg.Done()

	evChan := hook.Start()

	var lastX, lastY int
	var isDragging bool

	for {
		select {
		case <-c.stopChan:
			return
		case ev := <-evChan:
			if !c.isRunning {
				return
			}

			x, y := int(ev.X), int(ev.Y)
			currentTime := time.Now()

			// 根据事件类型分发处理
			switch {
			// 处理键盘事件
			case ev.Kind >= hook.KeyDown && ev.Kind <= hook.KeyHold:
				// 检查修饰键
				var modifiers []string
				if ev.Mask&0x0100 != 0 || ev.Mask&0x11 != 0 {
					modifiers = append(modifiers, "Ctrl")
				}
				if ev.Mask&0x0200 != 0 || ev.Mask&0x10 != 0 {
					modifiers = append(modifiers, "Shift")
				}
				if ev.Mask&0x0400 != 0 || ev.Mask&0x12 != 0 {
					modifiers = append(modifiers, "Alt")
				}
				if ev.Mask&0x0800 != 0 || ev.Mask&0x5B != 0 {
					modifiers = append(modifiers, "Win")
				}

				switch ev.Kind {
				case hook.KeyDown:
					c.keyboardChan <- model.KeyboardEvent{
						Type:      model.KeyDown,
						Keycode:   ev.Keychar,
						Keyname:   string(rune(ev.Keychar)),
						Timestamp: currentTime,
						Modifiers: modifiers,
					}
				case hook.KeyUp:
					c.keyboardChan <- model.KeyboardEvent{
						Type:      model.KeyUp,
						Keycode:   ev.Keychar,
						Keyname:   string(rune(ev.Keychar)),
						Timestamp: currentTime,
						Modifiers: modifiers,
					}
				case hook.KeyHold:
					c.keyboardChan <- model.KeyboardEvent{
						Type:      model.KeyPress,
						Keycode:   ev.Keychar,
						Keyname:   string(rune(ev.Keychar)),
						Timestamp: currentTime,
						Modifiers: modifiers,
					}
				}

			// 处理鼠标事件
			case ev.Kind >= hook.MouseDown && ev.Kind <= hook.MouseWheel:
				// 处理鼠标移动事件
				if x != lastX || y != lastY {
					if isDragging {
						c.mouseChan <- model.MouseEvent{
							Type:      model.MouseDrag,
							X:         x,
							Y:         y,
							Button:    "left",
							Timestamp: currentTime,
						}
					} else {
						c.mouseChan <- model.MouseEvent{
							Type:      model.MouseMove,
							X:         x,
							Y:         y,
							Timestamp: currentTime,
						}
					}
					lastX, lastY = x, y
				}

				// 处理鼠标按钮事件
				switch ev.Kind {
				case hook.MouseDown:
					button := "left"
					if ev.Button == 3 {
						button = "right"
					} else if ev.Button == 2 {
						button = "middle"
					}

					c.mouseChan <- model.MouseEvent{
						Type:      model.MouseClick,
						X:         x,
						Y:         y,
						Button:    button,
						Timestamp: currentTime,
					}

					isDragging = true

				case hook.MouseUp:
					isDragging = false

				case hook.MouseWheel:
					c.mouseChan <- model.MouseEvent{
						Type:      model.MouseScroll,
						X:         x,
						Y:         y,
						Delta:     int(ev.Amount),
						Timestamp: currentTime,
					}
				}
			}
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
