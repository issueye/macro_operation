package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	pb "github.com/issueye/macro-operation/macro-common/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	engineAddr = "localhost:50051"
)

// App 应用主结构
type App struct {
	ctx      context.Context
	conn     *grpc.ClientConn
	client   pb.MacroEngineClient
	engineCmd *exec.Cmd
}

// NewApp 创建新应用
func NewApp() *App {
	return &App{}
}

// startEngine 启动 engine 服务
func (a *App) startEngine() error {
	// 获取 engine 可执行文件路径
	exePath := filepath.Join(filepath.Dir(os.Args[0]), "engine.exe")
	if _, err := os.Stat(exePath); err != nil {
		// 尝试从当前目录查找
		exePath = "engine.exe"
	}

	log.Printf("Starting engine: %s", exePath)

	// 启动 engine 进程
	cmd := exec.Command(exePath)

	// 隐藏控制台窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	a.engineCmd = cmd
	log.Println("Engine process started")

	return nil
}

// waitForEngine 等待 engine 准备好
func (a *App) waitForEngine(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := grpc.Dial(engineAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second),
		)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("engine timeout")
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 启动 engine 服务
	if err := a.startEngine(); err != nil {
		log.Printf("Failed to start engine: %v", err)
		return
	}

	// 等待 engine 准备好
	if err := a.waitForEngine(10 * time.Second); err != nil {
		log.Printf("Failed to wait for engine: %v", err)
		return
	}

	// 连接到 gRPC 服务
	conn, err := grpc.Dial(engineAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return
	}

	a.conn = conn
	a.client = pb.NewMacroEngineClient(conn)

	log.Println("Connected to Macro Engine gRPC server")
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	// 关闭 gRPC 连接
	if a.conn != nil {
		a.conn.Close()
	}
	// 停止 engine 进程
	if a.engineCmd != nil && a.engineCmd.Process != nil {
		a.engineCmd.Process.Kill()
		a.engineCmd.Wait()
		log.Println("Engine process stopped")
	}
}

// IsEngineRunning 检查 engine 是否在运行
func (a *App) IsEngineRunning() bool {
	if a.engineCmd == nil || a.engineCmd.Process == nil {
		return false
	}
	// 检查进程是否存在
	process, err := os.FindProcess(a.engineCmd.Process.Pid)
	return err == nil && process.Signal(os.Signal(nil)) == nil
}

// ========== 录制相关 ==========

// StartRecording 开始录制
func (a *App) StartRecording() error {
	if a.client == nil {
		return fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.StartRecording(a.ctx, &pb.StartRecordingRequest{
		ClearPrevious: true,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

// StopRecording 停止录制
func (a *App) StopRecording() (int, error) {
	if a.client == nil {
		return 0, fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.StopRecording(a.ctx, &pb.StopRecordingRequest{})
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf(resp.Message)
	}

	return int(resp.EventCount), nil
}

// IsRecording 获取录制状态
func (a *App) IsRecording() bool {
	if a.client == nil {
		return false
	}

	resp, err := a.client.GetRecordingStatus(a.ctx, &pb.GetRecordingStatusRequest{})
	if err != nil {
		return false
	}

	return resp.IsRecording
}

// GetEventCount 获取事件数量
func (a *App) GetEventCount() int {
	if a.client == nil {
		return 0
	}

	resp, err := a.client.GetRecordingStatus(a.ctx, &pb.GetRecordingStatusRequest{})
	if err != nil {
		return 0
	}

	return int(resp.EventCount)
}

// ========== 回放相关 ==========

// PlayScript 播放脚本
func (a *App) PlayScript(script string) error {
	if a.client == nil {
		return fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.PlayScript(a.ctx, &pb.PlayScriptRequest{
		Script: script,
		Speed:  1.0,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

// ========== 宏管理相关 ==========

// SaveMacro 保存宏
func (a *App) SaveMacro(name, script string) error {
	if a.client == nil {
		return fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.SaveMacro(a.ctx, &pb.SaveMacroRequest{
		Name:   name,
		Script: script,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

// LoadMacro 加载宏
func (a *App) LoadMacro(name string) (string, error) {
	if a.client == nil {
		return "", fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.LoadMacro(a.ctx, &pb.LoadMacroRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf(resp.Message)
	}

	return resp.Macro.Script, nil
}

// DeleteMacro 删除宏
func (a *App) DeleteMacro(name string) error {
	if a.client == nil {
		return fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.DeleteMacro(a.ctx, &pb.DeleteMacroRequest{
		Name: name,
	})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

// MacroInfo 宏信息
type MacroInfo struct {
	Name      string
	CreatedAt string
	UpdatedAt string
	Script    string
}

// ListMacros 列出所有宏
func (a *App) ListMacros() ([]MacroInfo, error) {
	if a.client == nil {
		return nil, fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.ListMacros(a.ctx, &pb.ListMacrosRequest{})
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf(resp.Message)
	}

	result := make([]MacroInfo, 0, len(resp.Macros))
	for _, m := range resp.Macros {
		result = append(result, MacroInfo{
			Name:      m.Name,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			Script:    m.Script,
		})
	}

	return result, nil
}

// ========== 脚本生成相关 ==========

// GenerateScript 生成脚本
func (a *App) GenerateScript(name string) (string, error) {
	if a.client == nil {
		return "", fmt.Errorf("not connected to engine server")
	}

	resp, err := a.client.GenerateScript(a.ctx, &pb.GenerateScriptRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf(resp.Message)
	}

	return resp.Script, nil
}

// GenerateCurrentScript 生成当前录制事件的脚本
func (a *App) GenerateCurrentScript() (string, error) {
	return a.GenerateScript("temp")
}
