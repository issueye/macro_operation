package executor

import (
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/go-vgo/robotgo"
)

// Executor 脚本执行器
type Executor struct {
	isPlaying bool
	mutex     sync.RWMutex
	vm        *goja.Runtime
}

// NewExecutor 创建新的执行器
func NewExecutor() *Executor {
	return &Executor{
		isPlaying: false,
	}
}

// Execute 执行脚本
func (e *Executor) Execute(script string) error {
	e.mutex.Lock()
	if e.isPlaying {
		e.mutex.Unlock()
		return fmt.Errorf("playback already in progress")
	}
	e.isPlaying = true
	e.mutex.Unlock()

	defer func() {
		e.mutex.Lock()
		e.isPlaying = false
		e.mutex.Unlock()
	}()

	// 创建 JavaScript 运行时
	vm := goja.New()

	// 注册机器人 API
	if err := e.registerRobotAPI(vm); err != nil {
		return fmt.Errorf("failed to register robot API: %w", err)
	}

	// 执行脚本
	_, err := vm.RunString(script)
	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return nil
}

// IsPlaying 检查是否正在播放
func (e *Executor) IsPlaying() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.isPlaying
}

// registerRobotAPI 注册机器人 API
func (e *Executor) registerRobotAPI(vm *goja.Runtime) error {
	// 鼠标移动
	mouseMove := func(x, y int) {
		robotgo.MoveMouse(x, y)
	}

	// 鼠标点击
	mouseClick := func(button string) {
		robotgo.MouseClick(button)
	}

	// 鼠标拖拽
	mouseDrag := func(x, y int) {
		robotgo.DragMouse(x, y)
	}

	// 鼠标滚轮
	mouseScroll := func(delta int) {
		robotgo.ScrollDir(delta, "up")
	}

	// 键盘按下
	keyDown := func(key string) {
		robotgo.KeyDown(key)
	}

	// 键盘释放
	keyUp := func(key string) {
		robotgo.KeyUp(key)
	}

	// 键盘输入
	keyType := func(keys string) {
		robotgo.TypeStr(keys)
	}

	// 键盘敲击
	keyTap := func(key string) {
		robotgo.KeyTap(key)
	}

	// 睡眠
	sleep := func(ms int) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	// 截图
	screenshot := func() string {
		bitmap := robotgo.CaptureScreen(0, 0, 0, 0)
		defer robotgo.FreeBitmap(bitmap)
		img := robotgo.ToImage(bitmap)
		return robotgo.ToStringImg(img, "png")
	}

	// 日志
	log := func(msg string) {
		fmt.Println("[JS]", msg)
	}

	// 设置全局对象
	vm.Set("mouseMove", mouseMove)
	vm.Set("mouseClick", mouseClick)
	vm.Set("mouseDrag", mouseDrag)
	vm.Set("mouseScroll", mouseScroll)
	vm.Set("keyDown", keyDown)
	vm.Set("keyUp", keyUp)
	vm.Set("keyType", keyType)
	vm.Set("keyTap", keyTap)
	vm.Set("sleep", sleep)
	vm.Set("screenshot", screenshot)
	vm.Set("log", log)

	return nil
}
