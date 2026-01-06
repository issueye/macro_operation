//go:build windows
// +build windows

package selection

import (
	"syscall"
	"unsafe"
)

var (
	user32DLL     = syscall.MustLoadDLL("user32")
	kernel32DLL   = syscall.MustLoadDLL("kernel32")

	_GetForegroundWindow = user32DLL.MustFindProc("GetForegroundWindow")
	_GetWindowTextW     = user32DLL.MustFindProc("GetWindowTextW")
	_IsWindowVisible    = user32DLL.MustFindProc("IsWindowVisible")
	_OpenClipboard      = user32DLL.MustFindProc("OpenClipboard")
	_CloseClipboard     = user32DLL.MustFindProc("CloseClipboard")
	_EmptyClipboard     = user32DLL.MustFindProc("EmptyClipboard")
	_GetClipboardData   = user32DLL.MustFindProc("GetClipboardData")
	_SetClipboardData   = user32DLL.MustFindProc("SetClipboardData")
	_GlobalLock         = kernel32DLL.MustFindProc("GlobalLock")
	_GlobalUnlock       = kernel32DLL.MustFindProc("GlobalUnlock")
)

// GetSelectedText 获取当前选中的文本
// 通过模拟 Ctrl+C 复制并从剪贴板获取
func GetSelectedText() (string, error) {
	// 保存当前剪贴板内容
	originalText, err := GetClipboardText()
	if err != nil {
		// 如果获取失败，继续尝试
		originalText = ""
	}

	// 打开剪贴板
	r1, _, err := _OpenClipboard.Call(0)
	if r1 == 0 {
		return "", err
	}
	defer _CloseClipboard.Call()

	// 清空剪贴板
	_EmptyClipboard.Call()

	// 模拟 Ctrl+C
	SimulateCopy()

	// 等待复制完成
	text, err := GetClipboardText()
	if err != nil {
		return "", err
	}

	// 恢复原始剪贴板内容
	if originalText != "" {
		SetClipboardText(originalText)
	}

	return text, nil
}

// GetClipboardText 从剪贴板获取文本
func GetClipboardText() (string, error) {
	r1, _, err := _OpenClipboard.Call(0)
	if r1 == 0 {
		return "", err
	}
	defer _CloseClipboard.Call()

	// 获取 CF_UNICODETEXT 格式的数据
	handle, _, err := _GetClipboardData.Call(0xD) // CF_UNICODETEXT = 13
	if handle == 0 {
		return "", err
	}

	// 锁定全局内存
	ptr, _, err := _GlobalLock.Call(handle)
	if ptr == 0 {
		return "", err
	}
	defer _GlobalUnlock.Call(handle)

	// 转换为 Go 字符串
	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
	return text, nil
}

// SetClipboardText 设置剪贴板文本
func SetClipboardText(text string) error {
	r1, _, err := _OpenClipboard.Call(0)
	if r1 == 0 {
		return err
	}
	defer _CloseClipboard.Call()

	_EmptyClipboard.Call()

	data, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}

	// 计算所需内存大小（包含 null 终止符）
	size := len(data) * 2

	// 分配全局内存
	hMem, _, err := kernel32DLL.MustFindProc("GlobalAlloc").Call(0x0002, uintptr(size)) // GMEM_MOVEABLE = 0x0002
	if hMem == 0 {
		return err
	}

	// 锁定内存
	mem, _, err := _GlobalLock.Call(hMem)
	if mem == 0 {
		return err
	}
	defer _GlobalUnlock.Call(hMem)

	// 复制数据
	copy((*[1 << 20]byte)(unsafe.Pointer(mem))[:], unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), size))

	// 设置剪贴板数据
	_SetClipboardData.Call(0xD, hMem) // CF_UNICODETEXT = 13

	return nil
}

// SimulateCopy 模拟 Ctrl+C 复制操作
func SimulateCopy() {
	// 发送 Ctrl+C 按键
	// VK_CONTROL = 0x11
	// ord('C') = 0x43

	// keyDown Ctrl
	_KeyDown(0x11)
	// keyDown C
	_KeyDown(0x43)
	// keyUp C
	_KeyUp(0x43)
	// keyUp Ctrl
	_KeyUp(0x11)
}

func _KeyDown(vk byte) {
	// 使用 PostMessage 发送按键消息
	hwnd, _, _ := _GetForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	// WM_KEYDOWN = 0x0100
	user32DLL.MustFindProc("PostMessageW").Call(hwnd, 0x0100, uintptr(vk), 0)
}

func _KeyUp(vk byte) {
	hwnd, _, _ := _GetForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	// WM_KEYUP = 0x0101
	user32DLL.MustFindProc("PostMessageW").Call(hwnd, 0x0101, uintptr(vk), 0)
}
