package service

import (
	"fmt"
	"time"
)

// Example_recordServiceHotkeys 演示如何使用 RecordService 的内置快捷键
func Example_recordServiceHotkeys() {
	// 创建录制服务（会自动注册默认快捷键）
	recordService := NewRecordService()

	// 设置快捷键动作回调
	recordService.SetActionCallback(HotkeyActionToggleRecording, func() {
		fmt.Println("[回调] 录制状态已改变")
	})

	recordService.SetActionCallback(HotkeyActionSave, func() {
		script, err := recordService.GenerateScript("example")
		if err != nil {
			fmt.Printf("[回调] 生成脚本失败: %v\n", err)
			return
		}
		fmt.Printf("[回调] 脚本已生成，长度: %d 字符\n", len(script))
	})

	recordService.SetActionCallback(HotkeyActionClear, func() {
		fmt.Println("[回调] 录制已清空")
	})

	// 启动事件监听（底层会启动 hook.Start()）
	capture := recordService.GetCapture()
	capture.Start()

	fmt.Println("=== RecordService 内置快捷键 ===")
	fmt.Println("  Ctrl+R         - 开始/停止录制")
	fmt.Println("  Ctrl+S         - 保存")
	fmt.Println("  Ctrl+Shift+C   - 清空录制")
	fmt.Println("  Ctrl+Shift+R   - 停止录制并保存")
	fmt.Println()
	fmt.Println("现在可以使用快捷键控制录制了！")
	fmt.Println("程序将运行 30 秒后自动退出...")

	time.Sleep(30 * time.Second)

	capture.Stop()
	fmt.Println("程序已退出")
}

// Example_registerCustomHotkey 演示如何注册自定义快捷键
func Example_registerCustomHotkey() {
	recordService := NewRecordService()

	// 获取捕获器来注册自定义快捷键
	capture := recordService.GetCapture()

	// 注册播放快捷键
	capture.RegisterHotkey("ctrl+p", func() {
		fmt.Println("[快捷键] 播放宏")
		events := recordService.GetEvents()
		fmt.Printf("[快捷键] 准备播放 %d 个事件\n", len(events))
	})

	// 注册删除快捷键
	capture.RegisterHotkey("ctrl+d", func() {
		fmt.Println("[快捷键] 删除最后一个事件")
		events := recordService.GetEvents()
		if len(events) > 0 {
			// 注意：这里只是示例，实际实现可能需要修改 Capture 来支持删除事件
			fmt.Printf("[快捷键] 事件数量: %d\n", len(events))
		}
	})

	// 注册撤销快捷键
	capture.RegisterHotkey("ctrl+z", func() {
		fmt.Println("[快捷键] 撤销")
	})

	capture.Start()

	fmt.Println("=== 自定义快捷键已注册 ===")
	fmt.Println("  Ctrl+P - 播放宏")
	fmt.Println("  Ctrl+D - 删除")
	fmt.Println("  Ctrl+Z - 撤销")
	fmt.Println()

	time.Sleep(30 * time.Second)
	capture.Stop()
}

// Example_copyPasteHotkeys 演示复制粘贴快捷键
// 注意：这些快捷键会被正常记录到脚本中
func Example_copyPasteHotkeys() {
	recordService := NewRecordService()
	capture := recordService.GetCapture()

	// 这些快捷键会触发特殊动作，但也会被记录
	capture.RegisterHotkey("ctrl+c", func() {
		fmt.Println("[快捷键] 复制操作（将被记录到脚本）")
	})

	capture.RegisterHotkey("ctrl+v", func() {
		fmt.Println("[快捷键] 粘贴操作（将被记录到脚本）")
	})

	capture.RegisterHotkey("ctrl+x", func() {
		fmt.Println("[快捷键] 剪切操作（将被记录到脚本）")
	})

	capture.Start()

	fmt.Println("=== 复制粘贴快捷键 ===")
	fmt.Println("  Ctrl+C - 复制")
	fmt.Println("  Ctrl+V - 粘贴")
	fmt.Println("  Ctrl+X - 剪切")
	fmt.Println()
	fmt.Println("这些操作会被记录到脚本中，回放时会自动执行")

	time.Sleep(30 * time.Second)
	capture.Stop()
}
