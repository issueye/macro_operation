package model

// EventType 事件类型
type EventType string

const (
	EventTypeKeyDown   EventType = "keydown"
	EventTypeKeyUp     EventType = "keyup"
	EventTypeMouseDown EventType = "mousedown"
	EventTypeMouseUp   EventType = "mouseup"
	EventTypeMouseMove EventType = "mousemove"
	EventTypeWheel     EventType = "wheel"
)

// MouseButton 鼠标按钮
type MouseButton int

const (
	MouseButtonLeft   MouseButton = 1
	MouseButtonRight  MouseButton = 2
	MouseButtonMiddle MouseButton = 3
)

// Event 事件模型
type Event struct {
	Type      EventType  `json:"type"`
	KeyCode   int        `json:"key_code"`
	X         int        `json:"x"`
	Y         int        `json:"y"`
	Button    MouseButton `json:"button"`
	Timestamp int64      `json:"timestamp"`
	Delta     int        `json:"delta"` // 滚轮增量
}
