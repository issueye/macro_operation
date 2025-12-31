package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/issueye/macro_operation/configs"
	"github.com/issueye/macro_operation/internal/repository"
	"github.com/issueye/macro_operation/internal/service"
)

var (
	Version   = "1.0.0"
	BuildTime = "2024-12-31"
)

func main() {
	log.Printf("Starting Macro Recorder v%s (Build: %s)", Version, BuildTime)
	log.Println("操作宏录制和回放工具")

	// 加载配置
	configPath := configs.GetConfigPath()
	config, err := configs.Load(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load config from %s: %v", configPath, err)
		log.Println("Using default configuration")
		// 使用默认配置
		config = &configs.Config{}
		SetDefaults(config)
	}
	log.Printf("Config loaded: %+v", config.App)

	// 初始化服务
	macroRepo, err := repository.NewFileMacroRepository(config.Storage.MacrosDir)
	if err != nil {
		log.Fatalf("Failed to initialize macro repository: %v", err)
	}

	recordService := service.NewRecordService()
	playService := service.NewPlayService()
	macroService := service.NewMacroService(macroRepo)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nReceived interrupt signal, shutting down...")
		if recordService.IsStarted() {
			recordService.Stop()
		}
		os.Exit(0)
	}()

	// 运行CLI界面
	runCLI(recordService, playService, macroService)
}

// SetDefaults 设置默认配置
func SetDefaults(config *configs.Config) {
	if config.App.Name == "" {
		config.App.Name = "Macro Recorder"
	}
	if config.App.Version == "" {
		config.App.Version = "1.0.0"
	}
	if config.Storage.MacrosDir == "" {
		config.Storage.MacrosDir = "./macros"
	}
}

// runCLI 运行命令行界面
func runCLI(
	recordService *service.RecordService,
	playService *service.PlayService,
	macroService *service.MacroService,
) {
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("请选择操作 (1-7): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			handleRecord(reader, recordService)
		case "2":
			handleStop(recordService)
		case "3":
			HandleSave(reader, recordService, macroService)
		case "4":
			handlePlay(reader, playService, macroService)
		case "5":
			handleList(macroService)
		case "6":
			handleDelete(reader, macroService)
		case "7":
			handleGenerateAndSave(reader, recordService, macroService)
		case "0", "q", "Q":
			log.Println("Exiting...")
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}

		fmt.Println()
		fmt.Println("按 Enter 继续...")
		reader.ReadString('\n')
	}
}

// printMenu 打印菜单
func printMenu() {
	fmt.Println("\n========== 操作宏录制和回放 ==========")
	fmt.Println("1. 开始录制")
	fmt.Println("2. 停止录制")
	fmt.Println("3. 保存当前录制的宏")
	fmt.Println("4. 播放宏")
	fmt.Println("5. 列出所有宏")
	fmt.Println("6. 删除宏")
	fmt.Println("7. 录制并保存 (一键操作)")
	fmt.Println("0. 退出")
	fmt.Println("======================================")
}

// handleRecord 处理录制
func handleRecord(reader *bufio.Reader, recordService *service.RecordService) {
	if recordService.IsStarted() {
		fmt.Println("已经在录制中...")
		return
	}

	fmt.Println("开始录制... (按 Ctrl+C 停止)")
	if err := recordService.Start(); err != nil {
		fmt.Printf("启动录制失败: %v\n", err)
		return
	}
	fmt.Println("录制已开始，请执行您的操作")
	fmt.Println("录制中... (当前事件数: 0)")

	// 显示录制状态
	go func() {
		ticker := make(chan struct{})
		close(ticker)
		for range ticker {
			if !recordService.IsStarted() {
				return
			}
			count := recordService.GetEventCount()
			fmt.Printf("\r录制中... (当前事件数: %d) ", count)
		}
	}()

	fmt.Println("\n输入任何字符并按 Enter 停止录制...")
	reader.ReadString('\n')
}

// handleStop 处理停止
func handleStop(recordService *service.RecordService) {
	if !recordService.IsStarted() {
		fmt.Println("当前没有在录制")
		return
	}

	if err := recordService.Stop(); err != nil {
		fmt.Printf("停止录制失败: %v\n", err)
		return
	}

	eventCount := recordService.GetEventCount()
	fmt.Printf("录制已停止，共捕获 %d 个事件\n", eventCount)
}

// HandleSave 处理保存
func HandleSave(reader *bufio.Reader, recordService *service.RecordService, macroService *service.MacroService) {
	if recordService.IsStarted() {
		fmt.Println("请先停止录制")
		return
	}

	eventCount := recordService.GetEventCount()
	if eventCount == 0 {
		fmt.Println("没有录制的事件，无法保存")
		return
	}

	fmt.Print("请输入宏名称: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		fmt.Println("宏名称不能为空")
		return
	}

	script, err := recordService.GenerateScript(name)
	if err != nil {
		fmt.Printf("生成脚本失败: %v\n", err)
		return
	}

	if err := macroService.SaveFromScript(name, script); err != nil {
		fmt.Printf("保存宏失败: %v\n", err)
		return
	}

	fmt.Printf("宏 '%s' 已保存\n", name)
	fmt.Printf("\n生成的脚本:\n%s\n", script)
}

// handlePlay 处理播放
func handlePlay(reader *bufio.Reader, playService *service.PlayService, macroService *service.MacroService) {
	macros, err := macroService.List()
	if err != nil {
		fmt.Printf("获取宏列表失败: %v\n", err)
		return
	}

	if len(macros) == 0 {
		fmt.Println("没有可用的宏")
		return
	}

	fmt.Println("\n可用的宏:")
	for i, macro := range macros {
		fmt.Printf("%d. %s - %s\n", i+1, macro.Meta.Name, macro.Meta.Description)
	}

	fmt.Print("\n请输入要播放的宏名称或编号: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 尝试解析为编号
	if num, err := strconv.Atoi(input); err == nil {
		if num > 0 && num <= len(macros) {
			input = macros[num-1].Meta.Name
		}
	}

	macro, err := macroService.Load(input)
	if err != nil {
		fmt.Printf("加载宏失败: %v\n", err)
		return
	}

	fmt.Printf("正在播放宏 '%s'...\n", macro.Meta.Name)
	if err := playService.PlayScript(macro.Script.Code); err != nil {
		fmt.Printf("播放宏失败: %v\n", err)
		return
	}
	fmt.Println("播放完成")
}

// handleList 处理列表
func handleList(macroService *service.MacroService) {
	macros, err := macroService.List()
	if err != nil {
		fmt.Printf("获取宏列表失败: %v\n", err)
		return
	}

	if len(macros) == 0 {
		fmt.Println("没有保存的宏")
		return
	}

	fmt.Printf("\n共 %d 个宏:\n", len(macros))
	fmt.Println(strings.Repeat("-", 60))
	for i, macro := range macros {
		fmt.Printf("%d. 名称: %s\n", i+1, macro.Meta.Name)
		fmt.Printf("   描述: %s\n", macro.Meta.Description)
		fmt.Printf("   版本: %s\n", macro.Meta.Version)
		fmt.Printf("   创建时间: %s\n", macro.Meta.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println(strings.Repeat("-", 60))
	}
}

// handleDelete 处理删除
func handleDelete(reader *bufio.Reader, macroService *service.MacroService) {
	fmt.Print("请输入要删除的宏名称: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if !macroService.Exists(name) {
		fmt.Printf("宏 '%s' 不存在\n", name)
		return
	}

	fmt.Printf("确定要删除宏 '%s' 吗? (y/n): ", name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)

	if strings.ToLower(confirm) == "y" {
		if err := macroService.Delete(name); err != nil {
			fmt.Printf("删除宏失败: %v\n", err)
			return
		}
		fmt.Printf("宏 '%s' 已删除\n", name)
	} else {
		fmt.Println("已取消删除")
	}
}

// handleGenerateAndSave 一键录制并保存
func handleGenerateAndSave(reader *bufio.Reader, recordService *service.RecordService, macroService *service.MacroService) {
	fmt.Print("请输入宏名称: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		fmt.Println("宏名称不能为空")
		return
	}

	fmt.Println("开始录制... (按 Enter 停止)")
	if err := recordService.Start(); err != nil {
		fmt.Printf("启动录制失败: %v\n", err)
		return
	}

	fmt.Println("录制中，请执行您的操作...")
	reader.ReadString('\n')

	if err := recordService.Stop(); err != nil {
		fmt.Printf("停止录制失败: %v\n", err)
		return
	}

	eventCount := recordService.GetEventCount()
	fmt.Printf("录制已停止，共捕获 %d 个事件\n", eventCount)

	if eventCount == 0 {
		fmt.Println("没有录制的事件，无法保存")
		return
	}

	script, err := recordService.GenerateScript(name)
	if err != nil {
		fmt.Printf("生成脚本失败: %v\n", err)
		return
	}

	if err := macroService.SaveFromScript(name, script); err != nil {
		fmt.Printf("保存宏失败: %v\n", err)
		return
	}

	fmt.Printf("宏 '%s' 已保存\n", name)
	fmt.Printf("\n生成的脚本:\n%s\n", script)
}
