package events

import (
	"fmt"
	"sync"
	"time"

	"macro-engine/internal/model"
	"macro-engine/internal/selection"

	hook "github.com/robotn/gohook"
)

// Capture 事件捕获器
type Capture struct {
	events         []model.Event
	started        bool
	eventCh        chan model.Event
	stopCh         chan struct{}
	mutex          sync.RWMutex
	eventCallback  func(model.Event)     // 实时事件回调
	pendingKeyDown map[uint16]hook.Event // 待处理的 KeyDown 事件(keycode -> event)
	lastContent    string                // 上一次剪贴板内容
	isDragging     bool                  // 是否正在拖拽选择文字
	dragStartX     int                   // 拖拽起始 X 坐标
	dragStartY     int                   // 拖拽起始 Y 坐标
}

// NewCapture 创建新的捕获器
func NewCapture() *Capture {
	return &Capture{
		events:         make([]model.Event, 0),
		started:        false,
		eventCh:        make(chan model.Event, 1000),
		stopCh:         make(chan struct{}),
		pendingKeyDown: make(map[uint16]hook.Event),
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
	c.events = make([]model.Event, 0)
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
	c.events = make([]model.Event, 0)
	c.pendingKeyDown = make(map[uint16]hook.Event)
	c.isDragging = false
	c.dragStartX = 0
	c.dragStartY = 0
	c.mutex.Unlock()
}

// GetEvents 获取所有事件
func (c *Capture) GetEvents() []model.Event {
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

// SetEventCallback 设置实时事件回调
func (c *Capture) SetEventCallback(callback func(model.Event)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.eventCallback = callback
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
			// 处理鼠标按下 - 检测是否已有选中的文字
			if ev.Kind == hook.MouseDown && ev.Button == 1 {
				c.mutex.Lock()
				c.isDragging = true
				c.dragStartX = int(ev.X)
				c.dragStartY = int(ev.Y)
				c.mutex.Unlock()

				// 检测当前是否有选中的文字
				if text, err := selection.GetSelectedText(); err == nil && text != "" {
					selectionEvent := model.Event{
						Type:      model.Selection,
						X:         int(ev.X),
						Y:         int(ev.Y),
						Text:      text,
						Timestamp: time.Now().UnixMilli(),
					}
					c.mutex.Lock()
					c.events = append(c.events, selectionEvent)
					c.mutex.Unlock()

					if c.eventCallback != nil {
						c.eventCallback(selectionEvent)
					}
					fmt.Printf("检测到选中文字: %q\n", text)
				}
			}

			// 处理鼠标拖拽结束（MouseUp）
			if ev.Kind == hook.MouseUp && ev.Button == 1 {
				c.mutex.Lock()
				wasDragging := c.isDragging
				dragStartX, dragStartY := c.dragStartX, c.dragStartY
				c.isDragging = false
				c.mutex.Unlock()

				// 如果之前在拖拽，检测是否选中了文字
				if wasDragging {
					// 移动鼠标到当前坐标来获取选中的文字
					// 由于 gohook 已经捕获了鼠标移动，这里直接获取剪贴板内容
					if text, err := selection.GetSelectedText(); err == nil && text != "" {
						selectionEvent := model.Event{
							Type:      model.Selection,
							X:         int(ev.X),
							Y:         int(ev.Y),
							Text:      text,
							Timestamp: time.Now().UnixMilli(),
						}
						c.mutex.Lock()
						c.events = append(c.events, selectionEvent)
						c.mutex.Unlock()

						if c.eventCallback != nil {
							c.eventCallback(selectionEvent)
						}
						fmt.Printf("拖拽选择文字: %q\n", text)
					}
				}
				_ = dragStartX
				_ = dragStartY
			}

			// 处理键盘事件配对逻辑
			var modelEvent model.Event
			if ev.Kind == hook.KeyDown { // KeyDown
				// 如果 keycode 为 0,先记录下来,等待 KeyUp 事件来配对
				if ev.Keycode == 0 {
					// 使用特殊键 0 来标记,等待配对
					c.pendingKeyDown[0] = ev
					fmt.Printf("[DEBUG] Stored pending KeyDown with keycode=0, waiting for KeyUp\n")
					continue
				}
				// 正常的 KeyDown 事件,直接记录
				modelEvent = convertToModelEvent(ev)
			} else if ev.Kind == hook.KeyUp { // KeyUp
				// 检查是否有待处理的 keycode=0 的 KeyDown
				if pendingEv, ok := c.pendingKeyDown[0]; ok {
					// 找到配对,使用 KeyUp 的 keycode
					delete(c.pendingKeyDown, 0)
					// 修正 KeyDown 的 keycode 为 KeyUp 的 keycode
					pendingEv.Keycode = ev.Keycode
					modelEvent = convertToModelEvent(pendingEv)
					fmt.Printf("[DEBUG] Paired KeyDown(keycode=0) with KeyUp(keycode=%d)\n", ev.Keycode)

					// 记录修正后的 KeyDown
					c.mutex.Lock()
					c.events = append(c.events, modelEvent)
					c.mutex.Unlock()

					if c.eventCallback != nil {
						c.eventCallback(modelEvent)
					}
				}
				// 记录 KeyUp 事件
				modelEvent = convertToModelEvent(ev)
			} else if ev.Kind == hook.MouseDrag {
				// MouseDrag 事件直接转换
				modelEvent = convertToModelEvent(ev)
			} else if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseUp ||
				ev.Kind == hook.MouseMove || ev.Kind == hook.MouseWheel ||
				ev.Kind == hook.MouseHold {
				// 其他鼠标事件直接转换
				modelEvent = convertToModelEvent(ev)
			} else {
				// 其他事件类型
				modelEvent = convertToModelEvent(ev)
			}

			// 打印事件
			fmt.Printf("事件 -> %s\n", modelEvent.String())

			c.mutex.Lock()
			c.events = append(c.events, modelEvent)
			c.mutex.Unlock()

			// 调用实时事件回调
			if c.eventCallback != nil {
				c.eventCallback(modelEvent)
			}
		}
	}
}

// convertToModelEvent 将 gohook.Event 转换为 model.Event
func convertToModelEvent(ev hook.Event) model.Event {
	event := model.Event{
		Type:      getEventType(ev.Kind),
		KeyCode:   int(ev.Keycode),
		Chars:     int(ev.Rawcode),
		X:         int(ev.X),
		Y:         int(ev.Y),
		Button:    getMouseButton(int(ev.Button)),
		Timestamp: time.Now().UnixMilli(),
		Delta:     int(ev.Amount),
	}

	return event
}

// getEventType 获取事件类型
func getEventType(kind uint8) model.EventType {
	return model.EventType(kind)
}

// getMouseButton 获取鼠标按钮
func getMouseButton(btn int) model.MouseButton {
	switch btn {
	case 1:
		return model.MouseButtonLeft
	case 2:
		return model.MouseButtonRight
	case 3:
		return model.MouseButtonMiddle
	default:
		return model.MouseButtonLeft
	}
}
