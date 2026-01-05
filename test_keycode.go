package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func main() {
	// 监听单个按键
	evChan := hook.Start()
	defer hook.End()

	fmt.Println("请按下一些按键，我会显示它们的 keycode...")
	fmt.Println("按 ESC 退出")

	for ev := range evChan {
		if ev.Kind == 3 { // KeyDown
			if ev.Keycode == 1 { // ESC
				break
			}
			fmt.Printf("Keycode: %d, Keychar: %c (0x%X)\n", ev.Keycode, ev.Keychar, ev.Keychar)
		}
	}
}
