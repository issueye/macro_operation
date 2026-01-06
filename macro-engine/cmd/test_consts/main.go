package main

import (
	"fmt"
	hook "github.com/robotn/gohook"
)

func main() {
	fmt.Printf("KeyDown=%d, KeyUp=%d, KeyHold=%d\n", hook.KeyDown, hook.KeyUp, hook.KeyHold)
	fmt.Printf("MouseDown=%d, MouseUp=%d, MouseMove=%d\n", hook.MouseDown, hook.MouseUp, hook.MouseMove)
	fmt.Printf("MouseDrag=%d, MouseWheel=%d\n", hook.MouseDrag, hook.MouseWheel)
	fmt.Printf("MouseHold=%d\n", hook.MouseHold)
	fmt.Printf("CharUndefined=%d, WheelUp=%d, WheelDown=%d\n", hook.CharUndefined, hook.WheelUp, hook.WheelDown)
	fmt.Printf("HookEnabled=%d, HookDisabled=%d, FakeEvent=%d\n", hook.HookEnabled, hook.HookDisabled, hook.FakeEvent)
}
