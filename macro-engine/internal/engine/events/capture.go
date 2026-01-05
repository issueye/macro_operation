package events

import (
	"fmt"
	"sync"
	"time"

	"macro-engine/internal/model"

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
			// 过滤无效事件
			if !isValidEvent(ev) {
				fmt.Printf("[DEBUG] Filtered invalid event - Kind:%d, Keycode:%d\n", ev.Kind, ev.Keycode)
				continue
			}

			// 只记录 KeyDown 和 KeyUp 事件，忽略 KeyChar（字符输入事件）
			// 这样可以记录原始的按键序列，适用于任何输入法
			if ev.Kind == 4 {
				// 跳过 KeyChar 事件（Kind=4）
				// 调试：仍然打印以便了解情况
				fmt.Printf("[DEBUG] Skipped KeyChar event - Keycode:%d, Keychar:%d (0x%X)\n",
					ev.Keycode, ev.Keychar, ev.Keychar)
				continue
			}

			// 处理键盘事件配对逻辑
			var modelEvent model.Event
			if ev.Kind == 3 { // KeyDown
				// 如果 keycode 为 0,先记录下来,等待 KeyUp 事件来配对
				if ev.Keycode == 0 {
					// 使用特殊键 0 来标记,等待配对
					c.pendingKeyDown[0] = ev
					fmt.Printf("[DEBUG] Stored pending KeyDown with keycode=0, waiting for KeyUp\n")
					continue
				}
				// 正常的 KeyDown 事件,直接记录
				modelEvent = convertToModelEvent(ev)
			} else if ev.Kind == 5 { // KeyUp
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
			} else {
				// 其他事件类型
				modelEvent = convertToModelEvent(ev)
			}

			// 调试：打印记录的事件
			if ev.Kind >= 3 && ev.Kind <= 5 {
				fmt.Printf("[DEBUG] Recorded event - Kind:%d(%s), Keycode:%d (%s), Keychar:%d (0x%X)\n",
					ev.Kind, getEventKindName(ev.Kind), ev.Keycode, getKeyCodeName(ev.Keycode), ev.Keychar, ev.Keychar)
			}

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

// getEventKindName 获取事件类型名称
func getEventKindName(kind uint8) string {
	switch kind {
	case 3:
		return "KeyDown"
	case 4:
		return "KeyChar"
	case 5:
		return "KeyUp"
	default:
		return fmt.Sprintf("Unknown(%d)", kind)
	}
}

// getKeyCodeName 获取按键名称（用于调试）
func getKeyCodeName(keycode uint16) string {
	// 常见按键映射
	keyMap := map[uint16]string{
		1:  "left",   // 鼠标左键
		2:  "right",  // 鼠标右键
		4:  "middle", // 鼠标中键
		8:  "backspace",
		9:  "tab",
		13: "enter",
		16: "shift",
		17: "ctrl",
		18: "alt",
		27: "escape",
		32: "space",
		37: "left",
		38: "up",
		39: "right",
		40: "down",
		46: "delete",
		// 字母键
		65: "a", 66: "b", 67: "c", 68: "d", 69: "e", 70: "f", 71: "g", 72: "h",
		73: "i", 74: "j", 75: "k", 76: "l", 77: "m", 78: "n", 79: "o", 80: "p",
		81: "q", 82: "r", 83: "s", 84: "t", 85: "u", 86: "v", 87: "w", 88: "x",
		89: "y", 90: "z",
	}

	if name, ok := keyMap[keycode]; ok {
		return name
	}

	return fmt.Sprintf("key_%d", keycode)
}

// convertToModelEvent 将 gohook.Event 转换为 model.Event
func convertToModelEvent(ev hook.Event) model.Event {
	event := model.Event{
		Type:      getEventType(ev.Kind),
		KeyCode:   int(ev.Keycode),
		X:         int(ev.X),
		Y:         int(ev.Y),
		Button:    getMouseButton(int(ev.Button)),
		Timestamp: time.Now().UnixMilli(),
		Delta:     int(ev.Amount),
	}

	// 处理字符输入事件（KeyChar = 4）
	// Keychar 字段本身就是 rune 类型，可以直接转换为字符串
	// 它包含了实际输入的字符，包括中文等多字节字符
	if ev.Kind == 4 && ev.Keychar != 0 {
		// 验证是否为有效的 Unicode 字符（非替换字符）
		if ev.Keychar != 0xFFFD && ev.Keychar >= 32 && ev.Keychar <= 0x10FFFF {
			event.Chars = string(ev.Keychar)
		}
	}

	return event
}

// isValidEvent 检查事件是否有效
func isValidEvent(ev hook.Event) bool {
	// 过滤未知类型的事件（Kind 不在有效范围内）
	validKinds := map[uint8]bool{
		3:  true, // KeyDown
		4:  true, // KeyChar
		5:  true, // KeyUp
		6:  true, // MouseUp
		8:  true, // MouseDown
		9:  true, // MouseMove
		10: true, // MouseDrag
		11: true, // MouseWheel
	}

	if !validKinds[ev.Kind] {
		return false
	}

	// 对于键盘事件，放宽检查条件
	// KeyDown 事件：允许 keycode 为 0（某些输入法会产生这种情况）
	// KeyUp 事件：允许 keycode 为 0
	if ev.Kind == 3 || ev.Kind == 5 { // KeyDown 或 KeyUp
		// 不再过滤 keycode 为 0 的事件
		// 只过滤明显异常的大 keycode（大于 65535）
		if ev.Keycode > 65535 {
			return false
		}
	}

	return true
}

// getEventType 获取事件类型
func getEventType(kind uint8) model.EventType {
	switch kind {
	case 3: // KeyDown
		return model.EventTypeKeyDown
	case 4: // KeyChar - 字符输入
		return model.EventTypeChars
	case 5: // KeyUp
		return model.EventTypeKeyUp
	case 8: // MouseDown
		return model.EventTypeMouseDown
	case 6: // MouseUp
		return model.EventTypeMouseUp
	case 9: // MouseMove
		return model.EventTypeMouseMove
	case 11: // MouseWheel
		return model.EventTypeWheel
	default:
		return model.EventTypeKeyDown
	}
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
