package model

import (
	hook "github.com/robotn/gohook"
)

type EventType int

const (
	// HookEnabled  EventType = hook.HookEnabled
	// HookDisabled EventType = hook.HookDisabled

	// 键盘
	EventTypeKeyDown EventType = hook.KeyDown // = 4
	EventTypeKeyHold EventType = hook.KeyHold // = 3
	EventTypeKeyUp   EventType = hook.KeyUp   // = 5

	// 鼠标
	EventTypeMouseDown EventType = hook.MouseDown // = 7
	EventTypeMouseHold EventType = hook.MouseHold // = 8
	EventTypeMouseUp   EventType = hook.MouseUp   // = 6

	// 鼠标
	EventTypeMouseMove  EventType = hook.MouseMove  // = 9
	EventTypeMouseDrag  EventType = hook.MouseDrag  // = 10
	EventTypeMouseWheel EventType = hook.MouseWheel // = 11

	EventTypeFakeEvent EventType = hook.FakeEvent // = 12

	// Keychar could be v
	EventTypeCharUndefined EventType = hook.CharUndefined // = 65535
	EventTypeWheelUp       EventType = hook.WheelUp       // = -1
	EventTypeWheelDown     EventType = hook.WheelDown     // = 1

	// 保持向后兼容的别名
	KeyDown       = EventTypeKeyDown
	KeyHold       = EventTypeKeyHold
	KeyUp         = EventTypeKeyUp
	MouseDown     = EventTypeMouseDown
	MouseHold     = EventTypeMouseHold
	MouseUp       = EventTypeMouseUp
	MouseMove     = EventTypeMouseMove
	MouseDrag     = EventTypeMouseDrag
	MouseWheel    = EventTypeMouseWheel
	FakeEvent     = EventTypeFakeEvent
	CharUndefined = EventTypeCharUndefined
	WheelUp       = EventTypeWheelUp
	WheelDown     = EventTypeWheelDown
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
	Type      EventType   `json:"type"`
	KeyCode   int         `json:"key_code"`
	Chars     string      `json:"chars"` // 字符输入（用于中文等）
	X         int         `json:"x"`
	Y         int         `json:"y"`
	Button    MouseButton `json:"button"`
	Timestamp int64       `json:"timestamp"`
	Delta     int         `json:"delta"` // 滚轮增量
}

func GetKeyCode(key string) uint16 {
	return hook.KeychartoRawcode(key)
}

func GetKeyName(keyCode uint16) string {
	return hook.RawcodetoKeychar(keyCode)
}

var eventNameMap = map[EventType]string{
	KeyDown:       "KeyDown",
	KeyHold:       "KeyHold",
	KeyUp:         "KeyUp",
	MouseDown:     "MouseDown",
	MouseHold:     "MouseHold",
	MouseUp:       "MouseUp",
	MouseMove:     "MouseMove",
	MouseDrag:     "MouseDrag",
	MouseWheel:    "MouseWheel",
	FakeEvent:     "FakeEvent",
	CharUndefined: "CharUndefined",
	WheelUp:       "WheelUp",
	WheelDown:     "WheelDown",
}

var nameEventMap = map[string]EventType{
	"KeyDown":       KeyDown,
	"KeyHold":       KeyHold,
	"KeyUp":         KeyUp,
	"MouseDown":     MouseDown,
	"MouseHold":     MouseHold,
	"MouseUp":       MouseUp,
	"MouseMove":     MouseMove,
	"MouseDrag":     MouseDrag,
	"MouseWheel":    MouseWheel,
	"FakeEvent":     FakeEvent,
	"CharUndefined": CharUndefined,
	"WheelUp":       WheelUp,
	"WheelDown":     WheelDown,
}

func GetEventName(eventType EventType) string {
	return eventNameMap[eventType]
}

func GetEventType(eventName string) EventType {
	return nameEventMap[eventName]
}
