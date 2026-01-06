package generator

import (
	"fmt"
	"macro-engine/internal/model"
	"time"
)

// Example_sleepOptimization 展示 sleep 优化效果
func Example_sleepOptimization() {
	// 创建测试事件序列
	events := createTestEvents()

	// 生成脚本（使用默认 50ms 阈值）
	gen := NewGenerator()
	script1, _ := gen.Generate("test_default", events)
	fmt.Println("=== 默认阈值 (50ms) ===")
	fmt.Println(script1)

	// 生成脚本（使用 100ms 阈值，更少的 sleep）
	gen2 := NewGenerator(WithMinDelayThreshold(100))
	script2, _ := gen2.Generate("test_100ms", events)
	fmt.Println("=== 100ms 阈值（更少 sleep）===")
	fmt.Println(script2)

	// 生成脚本（使用 0ms 阈值，保留所有延迟）
	gen3 := NewGenerator(WithMinDelayThreshold(0))
	script3, _ := gen3.Generate("test_all_delays", events)
	fmt.Println("=== 0ms 阈值（保留所有延迟）===")
	fmt.Println(script3)
}

// createTestEvents 创建测试事件序列
func createTestEvents() []model.Event {
	now := time.Now().UnixMilli()
	events := []model.Event{
		// 快速打字 "hello" - 间隔 20ms
		{Type: model.EventTypeKeyDown, KeyCode: 35, Timestamp: now},      // h
		{Type: model.EventTypeKeyUp, KeyCode: 35, Timestamp: now + 20},
		{Type: model.EventTypeKeyDown, KeyCode: 38, Timestamp: now + 40},  // e
		{Type: model.EventTypeKeyUp, KeyCode: 38, Timestamp: now + 60},
		{Type: model.EventTypeKeyDown, KeyCode: 24, Timestamp: now + 80},  // l
		{Type: model.EventTypeKeyUp, KeyCode: 24, Timestamp: now + 100},
		{Type: model.EventTypeKeyDown, KeyCode: 24, Timestamp: now + 120}, // l
		{Type: model.EventTypeKeyUp, KeyCode: 24, Timestamp: now + 140},
		{Type: model.EventTypeKeyDown, KeyCode: 31, Timestamp: now + 160}, // o
		{Type: model.EventTypeKeyUp, KeyCode: 31, Timestamp: now + 180},

		// 暂停 200ms
		{Type: model.EventTypeKeyDown, KeyCode: 28, Timestamp: now + 380}, // enter

		// 鼠标移动 - 间隔 10ms
		{Type: model.EventTypeMouseMove, X: 100, Y: 100, Timestamp: now + 400},
		{Type: model.EventTypeMouseMove, X: 110, Y: 110, Timestamp: now + 410},
		{Type: model.EventTypeMouseMove, X: 120, Y: 120, Timestamp: now + 420},
		{Type: model.EventTypeMouseDown, Button: 1, Timestamp: now + 500},
		{Type: model.EventTypeMouseUp, Button: 1, Timestamp: now + 520},
	}
	return events
}

// Example_sleepOptimizationComparison 比较 sleep 数量
func Example_sleepOptimizationComparison() {
	events := createTestEvents()

	// 统计不同阈值下的 sleep 数量
	thresholds := []int{0, 10, 50, 100, 200}
	fmt.Println("=== Sleep 数量对比 ===")
	for _, threshold := range thresholds {
		gen := NewGenerator(WithMinDelayThreshold(threshold))
		script, _ := gen.Generate("test", events)

		// 统计 sleep 行数
		sleepCount := countSleepCalls(script)
		fmt.Printf("阈值 %3dms: %2d 个 sleep 语句\n", threshold, sleepCount)
	}
}

func countSleepCalls(script string) int {
	count := 0
	for i := 0; i < len(script); i++ {
		if i+5 < len(script) && script[i:i+5] == "sleep" {
			count++
		}
	}
	return count
}

// Example_customThreshold 使用自定义阈值
func Example_customThreshold() {
	events := createTestEvents()

	// 创建 150ms 阈值的生成器
	gen := NewGenerator(WithMinDelayThreshold(150))
	script, _ := gen.Generate("custom_threshold", events)

	fmt.Println("=== 自定义阈值 150ms ===")
	fmt.Println(script)
}
