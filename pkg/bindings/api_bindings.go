package bindings

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

// APIBindings API绑定
type APIBindings struct{}

// NewAPIBindings 创建API绑定
func NewAPIBindings() *APIBindings {
	return &APIBindings{}
}

// MouseMove 鼠标移动
func (b *APIBindings) MouseMove(x, y int) error {
	robotgo.MoveMouse(x, y)
	return nil
}

// MouseClick 鼠标点击
func (b *APIBindings) MouseClick(button string) error {
	switch button {
	case "left":
		robotgo.MouseClick("left")
	case "right":
		robotgo.MouseClick("right")
	case "middle", "center":
		robotgo.MouseClick("center")
	default:
		return fmt.Errorf("unsupported mouse button: %s", button)
	}
	return nil
}

// MouseDrag 鼠标拖拽
func (b *APIBindings) MouseDrag(x, y int) error {
	robotgo.DragMouse(x, y, "left")
	return nil
}

// MouseScroll 鼠标滚动
func (b *APIBindings) MouseScroll(delta int) error {
	// robotgo.Scroll需要两个参数：x和y方向的滚动量
	robotgo.Scroll(0, delta)
	return nil
}

// KeyDown 按键按下
func (b *APIBindings) KeyDown(key string) error {
	return robotgo.KeyToggle(key, "down")
}

// KeyUp 按键释放
func (b *APIBindings) KeyUp(key string) error {
	return robotgo.KeyToggle(key, "up")
}

// KeyType 输入文本
func (b *APIBindings) KeyType(text string) error {
	robotgo.TypeStr(text)
	return nil
}

// KeyTap 按键点击
func (b *APIBindings) KeyTap(key string) error {
	return robotgo.KeyTap(key)
}

// Sleep 延迟等待
func (b *APIBindings) Sleep(ms int) error {
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}

// Screenshot 屏幕截图
func (b *APIBindings) Screenshot(path string) error {
	// TODO: 截图功能待完善
	// robotgo 的 API 在不断变化，暂时返回未实现错误
	return fmt.Errorf("screenshot feature not yet implemented")
}

// GetMousePos 获取鼠标位置
func (b *APIBindings) GetMousePos() (x, y int) {
	return robotgo.GetMousePos()
}

// GetScreenSize 获取屏幕尺寸
func (b *APIBindings) GetScreenSize() (width, height int) {
	return robotgo.GetScreenSize()
}
