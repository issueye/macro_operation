package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"macro-engine/internal/engine/events"
	"macro-engine/internal/model"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("键盘事件测试程序")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("说明：")
	fmt.Println("1. 程序启动后，请切换到记事本")
	fmt.Println("2. 在记事本中输入一些字符（如 'test'）")
	fmt.Println("3. 然后按 Ctrl+C 停止")
	fmt.Println("4. 查看捕获到的事件")
	fmt.Println()

	capture := events.NewCapture()

	if err := capture.Start(); err != nil {
		fmt.Printf("启动失败: %v\n", err)
		os.Exit(1)
	}
	defer capture.Stop()

	fmt.Println("开始录制... (按 Ctrl+C 停止)")
	fmt.Println()

	// 设置信号处理
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 每 2 秒打印一次状态
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-sigCh:
			fmt.Println("\n停止录制...")
			break loop
		case <-ticker.C:
			count := capture.GetEventCount()
			fmt.Printf("\r已捕获 %d 个事件...", count)
		}
	}

	// 获取所有事件
	recordedEvents := capture.GetEvents()
	fmt.Printf("\n\n总共捕获了 %d 个事件\n\n", len(recordedEvents))

	// 统计事件类型
	keyDownCount := 0
	keyUpCount := 0
	mouseMoveCount := 0
	mouseClickCount := 0

	for i, ev := range recordedEvents {
		switch ev.Type {
		case model.EventTypeKeyDown:
			keyDownCount++
			if i < 20 {
				fmt.Printf("[%-3d] KeyDown   - KeyCode: %d\n", i, ev.KeyCode)
			}
		case model.EventTypeKeyUp:
			keyUpCount++
			if i < 20 {
				fmt.Printf("[%-3d] KeyUp     - KeyCode: %d\n", i, ev.KeyCode)
			}
		case model.EventTypeMouseMove:
			mouseMoveCount++
		case model.EventTypeMouseDown:
			mouseClickCount++
		case model.EventTypeMouseUp:
			mouseClickCount++
		}
	}

	fmt.Println()
	fmt.Println("========== 事件统计 ==========")
	fmt.Printf("KeyDown:  %d\n", keyDownCount)
	fmt.Printf("KeyUp:    %d\n", keyUpCount)
	fmt.Printf("MouseMove: %d\n", mouseMoveCount)
	fmt.Printf("MouseClick:%d\n", mouseClickCount)
	fmt.Println("==============================")
	fmt.Println()

	if keyDownCount == 0 && keyUpCount == 0 {
		fmt.Println("⚠️  警告：没有捕获到键盘事件！")
		fmt.Println()
		fmt.Println("可能的原因：")
		fmt.Println("1. 你在应用内的文本框中输入，而不是在记事本中输入")
		fmt.Println("2. 输入法的问题（某些输入法会阻止键盘钩子）")
		fmt.Println("3. 权限不足（需要管理员权限）")
		fmt.Println()
		fmt.Println("建议：")
		fmt.Println("- 请切换到记事本或其他程序进行输入")
		fmt.Println("- 尝试使用英文输入法")
		fmt.Println("- 以管理员身份运行程序")
	} else {
		fmt.Println("✓ 成功捕获到键盘事件")
	}
}
