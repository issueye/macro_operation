package model

import (
	hook "github.com/robotn/gohook"
)

type EventType int

const (
	// HookEnabled  EventType = hook.HookEnabled
	// HookDisabled EventType = hook.HookDisabled

	// 键盘
	KeyDown EventType = hook.KeyDown
	KeyHold EventType = hook.KeyHold
	KeyUp   EventType = hook.KeyUp

	// 鼠标
	MouseDown EventType = hook.MouseDown
	MouseHold EventType = hook.MouseHold
	MouseUp   EventType = hook.MouseUp

	// 鼠标
	MouseMove  EventType = hook.MouseMove
	MouseDrag  EventType = hook.MouseDrag
	MouseWheel EventType = hook.MouseWheel

	FakeEvent EventType = hook.FakeEvent

	// Keychar could be v
	CharUndefined EventType = hook.CharUndefined
	WheelUp       EventType = hook.WheelUp
	WheelDown     EventType = hook.WheelDown
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
