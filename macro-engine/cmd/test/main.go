package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"macro-engine/internal/engine/events"
)

func main() {
	fmt.Println("Macro Engine Test Program")
	fmt.Println("========================")
	fmt.Println("This program tests keyboard/mouse event capture")
	fmt.Println("Press Ctrl+C to stop\n")

	// 创建事件捕获器
	capture := events.NewCapture()

	// 开始捕获
	if err := capture.Start(); err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}
	defer capture.Stop()

	fmt.Println("Recording started...")
	fmt.Println("Press some keys or move your mouse...")
	fmt.Println()

	// 设置信号处理
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 定时打印事件统计
	ticker := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker:
				count := capture.GetEventCount()
				fmt.Printf("\rEvents captured: %d", count)
			}
		}
	}()

	// 等待中断信号
	<-sigCh
	fmt.Println("\n\nStopping...")

	// 获取所有事件
	recordedEvents := capture.GetEvents()
	fmt.Printf("Total events captured: %d\n", len(recordedEvents))

	// 打印前 20 个事件详情
	fmt.Println("\nFirst 20 events:")
	for i, ev := range recordedEvents {
		if i >= 20 {
			break
		}
		fmt.Printf("%d. Type:%s, KeyCode:%d, Chars:%q, X:%d, Y:%d\n",
			i+1, ev.Type, ev.KeyCode, ev.Chars, ev.X, ev.Y)
	}

	fmt.Println("\nTest completed.")
}
