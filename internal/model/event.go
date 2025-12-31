package model

import "time"

// EventType 事件类型
type EventType int

const (
	KeyDown EventType = iota
	KeyUp
	KeyPress
	MouseMove
	MouseClick
	MouseDrag
	MouseScroll
)

// Event 事件接口
type Event interface {
	GetTimestamp() time.Time
	GetType() EventType
	String() string
}

// KeyboardEvent 键盘事件
type KeyboardEvent struct {
	Type      EventType
	Keycode   rune
	Keyname   string
	Timestamp time.Time
	Modifiers []string // 修饰键列表，如 ["Ctrl", "Shift"]
}

func (ke KeyboardEvent) GetTimestamp() time.Time {
	return ke.Timestamp
}

func (ke KeyboardEvent) GetType() EventType {
	return ke.Type
}

func (ke KeyboardEvent) String() string {
	action := "KeyDown"
	if ke.Type == KeyUp {
		action = "KeyUp"
	} else if ke.Type == KeyPress {
		action = "KeyPress"
	}
	return action + "(" + string(ke.Keycode) + ")"
}

// MouseEvent 鼠标事件
type MouseEvent struct {
	Type      EventType
	X, Y      int
	Button    string
	Delta     int
	Timestamp time.Time
}

func (me MouseEvent) GetTimestamp() time.Time {
	return me.Timestamp
}

func (me MouseEvent) GetType() EventType {
	return me.Type
}

func (me MouseEvent) String() string {
	switch me.Type {
	case MouseMove:
		return "MouseMove"
	case MouseClick:
		return "MouseClick(" + me.Button + ")"
	case MouseDrag:
		return "MouseDrag"
	case MouseScroll:
		return "MouseScroll(" + string(rune(me.Delta)) + ")"
	default:
		return "Unknown"
	}
}
