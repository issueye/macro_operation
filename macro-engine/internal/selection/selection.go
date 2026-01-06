//go:build !windows
// +build !windows

package selection

// GetSelectedText 获取当前选中的文本（非 Windows 平台返回空）
func GetSelectedText() (string, error) {
	return "", nil
}

// GetClipboardText 从剪贴板获取文本（非 Windows 平台返回空）
func GetClipboardText() (string, error) {
	return "", nil
}

// SetClipboardText 设置剪贴板文本（非 Windows 平台为空操作）
func SetClipboardText(text string) error {
	return nil
}

// SimulateCopy 模拟 Ctrl+C 复制操作（非 Windows 平台为空操作）
func SimulateCopy() {
}
