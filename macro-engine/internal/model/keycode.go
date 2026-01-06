package model

import "fmt"

// KeyCodeMapping 统一的键码映射
// gohook 返回的是扫描码 (Scan Code)，参考: https://www.win.tue.nl/~aeb/linux/kbd/scancodes-1.html
type KeyCodeMapping struct {
	ScanCode int    // 扫描码
	Name     string // robotgo 识别的键名
	Alias    string // 别名（可选）
}

// KeyCodeMaps 键码映射表（按扫描码组织）
var KeyCodeMaps = map[int]string{
	// === 特殊键 (01-3F) ===
	1:  "escape",      // 01 - ESC
	2:  "1",           // 02 - 1
	3:  "2",           // 03 - 2
	4:  "3",           // 04 - 3
	5:  "4",           // 05 - 4
	6:  "5",           // 06 - 5
	7:  "6",           // 07 - 6
	8:  "7",           // 08 - 7
	9:  "8",           // 09 - 8
	10: "9",           // 0A - 9
	11: "0",           // 0B - 0
	12: "-",           // 0C - -
	13: "=",           // 0D - =
	14: "backspace",   // 0E - Backspace
	15: "tab",         // 0F - Tab
	16: "q",           // 10 - Q
	17: "w",           // 11 - W
	18: "e",           // 12 - E
	19: "r",           // 13 - R
	20: "t",           // 14 - T
	21: "y",           // 15 - Y
	22: "u",           // 16 - U
	23: "i",           // 17 - I
	24: "o",           // 18 - O
	25: "p",           // 19 - P
	26: "[",           // 1A - [
	27: "]",           // 1B - ]
	28: "return",      // 1C - Enter (robotgo 使用 "return")
	29: "lctrl",       // 1D - Left Ctrl
	30: "a",           // 1E - A
	31: "s",           // 1F - S
	32: "d",           // 20 - D
	33: "f",           // 21 - F
	34: "g",           // 22 - G
	35: "h",           // 23 - H
	36: "j",           // 24 - J
	37: "k",           // 25 - K
	38: "l",           // 26 - L
	39: ";",           // 27 - ;
	40: "'",           // 28 - '
	41: "`",           // 29 - `
	42: "lshift",      // 2A - Left Shift
	43: "\\",          // 2B - \
	44: "z",           // 2C - Z
	45: "x",           // 2D - X
	46: "c",           // 2E - C
	47: "v",           // 2F - V
	48: "b",           // 30 - B
	49: "n",           // 31 - N
	50: "m",           // 32 - M
	51: ",",           // 33 - ,
	52: ".",           // 34 - .
	53: "/",           // 35 - /
	54: "rshift",      // 36 - Right Shift
	55: "kp_multiply", // 37 - Numpad *
	56: "lalt",        // 38 - Left Alt
	57: "space",       // 39 - Space
	58: "capslock",    // 3A - Caps Lock
	59: "f1",          // 3B - F1
	60: "f2",          // 3C - F2
	61: "f3",          // 3D - F3
	62: "f4",          // 3E - F4
	63: "f5",          // 3F - F5
	64: "f6",          // 40 - F6
	65: "f7",          // 41 - F7
	66: "f8",          // 42 - F8
	67: "f9",          // 43 - F9
	68: "f10",         // 44 - F10
	69: "numlock",     // 45 - Num Lock
	70: "scrolllock",  // 46 - Scroll Lock
	87: "f11",         // 57 - F11
	88: "f12",         // 58 - F12

	// === 数字键盘 (47-52) ===
	71: "kp_7",    // 47 - Numpad 7
	72: "kp_8",    // 48 - Numpad 8
	73: "kp_9",    // 49 - Numpad 9
	74: "kp_minus", // 4A - Numpad -
	75: "kp_4",    // 4B - Numpad 4
	76: "kp_5",    // 4C - Numpad 5
	77: "kp_6",    // 4D - Numpad 6
	78: "kp_plus", // 4E - Numpad +
	79: "kp_1",    // 4F - Numpad 1
	80: "kp_2",    // 50 - Numpad 2
	81: "kp_3",    // 51 - Numpad 3
	82: "kp_0",    // 52 - Numpad 0
	83: "kp_decimal", // 53 - Numpad .

	// === 扩展扫描码 (E0 前缀，扫描码 110+) ===
	// 方向键和编辑键
	110: "home",       // E0 47 - Home
	111: "up",         // E0 48 - Up
	112: "pageup",     // E0 49 - Page Up
	113: "left",       // E0 4B - Left
	114: "right",      // E0 4D - Right
	115: "end",        // E0 4F - End
	116: "down",       // E0 50 - Down
	117: "pagedown",   // E0 51 - Page Down
	118: "insert",     // E0 52 - Insert
	119: "delete",     // E0 53 - Delete
	125: "lwin",       // E0 5B - Left Windows
	126: "rwin",       // E0 5C - Right Windows
	127: "menu",       // E0 5D - Menu (Application key)
	138: "rctrl",      // E0 1D - Right Ctrl
	184: "ralt",       // E0 38 - Right Alt (AltGr)
	186: ";",          // E0 27 - 也有一些键盘映射
	187: "=",          // E0 2D - 也有一些键盘映射
	188: ",",          // E0 33 - 也有一些键盘映射
	189: "-",          // E0 35 - 也有一些键盘映射
	190: ".",          // E0 34 - 也有一些键盘映射
	191: "/",          // E0 2B - 也有一些键盘映射
	192: "`",          // E0 29 - 也有一些键盘映射

	// === 媒体键和其他扩展键 ===
	150: "sleep",      // E0 5E - Sleep
	151: "wake",       // E0 63 - Wake
	152: "power",      // E0 5F - Power

	// === 浏览器多媒体键 ===
	173: "mute",       // E0 20 - Mute
	174: "volumedown", // E0 2E - Volume Down
	175: "volumeup",   // E0 30 - Volume Up
	176: "nexttrack",  // E0 19 - Next Track
	177: "prevtrack",  // E0 10 - Previous Track
	178: "stopcd",     // E0 24 - Stop CD
	179: "playpause",  // E0 22 - Play/Pause

	// === 浏览器键 ===
	180: "mail",       // E0 6C - Mail
	181: "mediaselect",// E0 6D - Media Select
	182: "webhome",    // E0 32 - Browser Home
	183: "refresh",    // E0 21 - Refresh
	// 移除重复的映射 184, 186
	193: "f13",        // E0 68 - F13
	194: "f14",        // E0 69 - F14
	195: "f15",        // E0 6A - F15

	// === 日文键盘特殊键 ===
	// 如果需要支持日文键盘，可以添加更多映射
}

// GetKeyNameByScanCode 根据扫描码获取键名
func GetKeyNameByScanCode(scanCode int) string {
	if name, ok := KeyCodeMaps[scanCode]; ok {
		return name
	}
	return fmt.Sprintf("key_%d", scanCode)
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
	"enter":       "return",
	"backspace":   "backspace",

	// 数字键盘
	"num_lock":    "numlock",
	"scroll_lock": "scrolllock",
	"caps_lock":   "capslock",

	// 方向键
	"page_up":     "pageup",
	"page_down":   "pagedown",
}

// NormalizeKeyName 规范化键名（转换为 robotgo 格式）
func NormalizeKeyName(name string) string {
	if normalized, ok := RobotgoKeyName[name]; ok {
		return normalized
	}
	return name
}
