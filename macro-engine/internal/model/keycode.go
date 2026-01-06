package model

import (
	hook "github.com/robotn/gohook"
)

// KeyCodeMapping 统一的键码映射
// gohook 返回的是扫描码 (Scan Code)，参考: https://www.win.tue.nl/~aeb/linux/kbd/scancodes-1.html
type KeyCodeMapping struct {
	ScanCode int    // 扫描码
	Name     string // robotgo 识别的键名
	Alias    string // 别名（可选）
}

// GetKeyNameByScanCode 根据扫描码获取键名
func GetKeyNameByScanCode(scanCode int) string {
	return hook.RawcodetoKeychar(uint16(scanCode))
}

// RobotgoKeyName 兼容 robotgo 的键名映射
// robotgo 使用特定的键名格式，这里提供转换
var RobotgoKeyName = map[string]string{
	// 修饰键
	"left ctrl":   "lctrl",
	"right ctrl":  "rctrl",
	"left shift":  "lshift",
	"right shift": "rshift",
	"left alt":    "lalt",
	"right alt":   "ralt",
	"left win":    "lwin",
	"right win":   "rwin",

	// 特殊键
	"enter":     "return",
	"backspace": "backspace",

	// 数字键盘
	"num_lock":    "numlock",
	"scroll_lock": "scrolllock",
	"caps_lock":   "capslock",

	// 方向键
	"page_up":   "pageup",
	"page_down": "pagedown",
}

// NormalizeKeyName 规范化键名（转换为 robotgo 格式）
func NormalizeKeyName(name string) string {
	if normalized, ok := RobotgoKeyName[name]; ok {
		return normalized
	}
	return name
}
