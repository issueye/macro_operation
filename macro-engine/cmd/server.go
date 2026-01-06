package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"macro-engine/internal/model"
	"macro-engine/internal/service"

	pb "github.com/issueye/macro-operation/macro-common/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	port   = flag.Int("port", 50051, "gRPC server port")
	macros = flag.String("macros", "./macros", "macros directory")
)

// MacroEngineServer gRPC 服务实现
type MacroEngineServer struct {
	pb.UnimplementedMacroEngineServer
	recordService *service.RecordService
	playService   *service.PlayService
	macroService  *service.MacroService
}

// NewMacroEngineServer 创建新的服务
func NewMacroEngineServer() (*MacroEngineServer, error) {
	repo, err := service.NewFileMacroRepository(*macros)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &MacroEngineServer{
		recordService: service.NewRecordService(),
		playService:   service.NewPlayService(),
		macroService:  service.NewMacroService(repo),
	}, nil
}

// StartRecording 开始录制
func (s *MacroEngineServer) StartRecording(ctx context.Context, req *pb.StartRecordingRequest) (*pb.StartRecordingResponse, error) {
	err := s.recordService.Start(req.ClearPrevious)
	if err != nil {
		return &pb.StartRecordingResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.StartRecordingResponse{
		Success: true,
		Message: "recording started",
	}, nil
}

// StopRecording 停止录制
func (s *MacroEngineServer) StopRecording(ctx context.Context, req *pb.StopRecordingRequest) (*pb.StopRecordingResponse, error) {
	count, err := s.recordService.Stop()
	if err != nil {
		return &pb.StopRecordingResponse{
			Success:    false,
			EventCount: 0,
			Message:    err.Error(),
		}, nil
	}

	return &pb.StopRecordingResponse{
		Success:    true,
		EventCount: int32(count),
		Message:    "recording stopped",
	}, nil
}

// GetRecordingStatus 获取录制状态
func (s *MacroEngineServer) GetRecordingStatus(ctx context.Context, req *pb.GetRecordingStatusRequest) (*pb.GetRecordingStatusResponse, error) {
	return &pb.GetRecordingStatusResponse{
		IsRecording: s.recordService.IsStarted(),
		EventCount:  int32(s.recordService.GetEventCount()),
	}, nil
}

// PlayScript 播放脚本
func (s *MacroEngineServer) PlayScript(ctx context.Context, req *pb.PlayScriptRequest) (*pb.PlayScriptResponse, error) {
	speed := req.Speed
	if speed <= 0 {
		speed = 1.0
	}

	err := s.playService.PlayScript(req.Script, speed)
	if err != nil {
		return &pb.PlayScriptResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.PlayScriptResponse{
		Success: true,
		Message: "playback completed",
	}, nil
}

// GetPlaybackStatus 获取播放状态
func (s *MacroEngineServer) GetPlaybackStatus(ctx context.Context, req *pb.GetPlaybackStatusRequest) (*pb.GetPlaybackStatusResponse, error) {
	return &pb.GetPlaybackStatusResponse{
		IsPlaying: s.playService.IsPlaying(),
		Progress:  100, // TODO: 实现进度跟踪
	}, nil
}

// SaveMacro 保存宏
func (s *MacroEngineServer) SaveMacro(ctx context.Context, req *pb.SaveMacroRequest) (*pb.SaveMacroResponse, error) {
	err := s.macroService.SaveFromScript(req.Name, req.Script)
	if err != nil {
		return &pb.SaveMacroResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SaveMacroResponse{
		Success: true,
		Message: "macro saved",
	}, nil
}

// LoadMacro 加载宏
func (s *MacroEngineServer) LoadMacro(ctx context.Context, req *pb.LoadMacroRequest) (*pb.LoadMacroResponse, error) {
	macro, err := s.macroService.Load(req.Name)
	if err != nil {
		return &pb.LoadMacroResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoadMacroResponse{
		Success: true,
		Macro:   convertMacroToProto(macro),
		Message: "macro loaded",
	}, nil
}

// DeleteMacro 删除宏
func (s *MacroEngineServer) DeleteMacro(ctx context.Context, req *pb.DeleteMacroRequest) (*pb.DeleteMacroResponse, error) {
	err := s.macroService.Delete(req.Name)
	if err != nil {
		return &pb.DeleteMacroResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteMacroResponse{
		Success: true,
		Message: "macro deleted",
	}, nil
}

// ListMacros 列出所有宏
func (s *MacroEngineServer) ListMacros(ctx context.Context, req *pb.ListMacrosRequest) (*pb.ListMacrosResponse, error) {
	macros, err := s.macroService.List()
	if err != nil {
		return &pb.ListMacrosResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoMacros := make([]*pb.MacroInfo, 0, len(macros))
	for _, m := range macros {
		protoMacros = append(protoMacros, convertMacroToProto(m))
	}

	return &pb.ListMacrosResponse{
		Success: true,
		Macros:  protoMacros,
		Message: "success",
	}, nil
}

// GenerateScript 生成脚本
func (s *MacroEngineServer) GenerateScript(ctx context.Context, req *pb.GenerateScriptRequest) (*pb.GenerateScriptResponse, error) {
	script, err := s.recordService.GenerateScript(req.Name)
	if err != nil {
		return &pb.GenerateScriptResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GenerateScriptResponse{
		Success: true,
		Script:  script,
		Message: "script generated",
	}, nil
}

// GetCurrentEvents 获取当前事件
func (s *MacroEngineServer) GetCurrentEvents(ctx context.Context, req *pb.GetCurrentEventsRequest) (*pb.GetCurrentEventsResponse, error) {
	events := s.recordService.GetEvents()

	pbEvents := make([]*pb.EventData, 0, len(events))
	for _, ev := range events {
		pbEvents = append(pbEvents, &pb.EventData{
			Type:      string(ev.Type),
			KeyCode:   int32(ev.KeyCode),
			X:         int32(ev.X),
			Y:         int32(ev.Y),
			Button:    int32(ev.Button),
			Timestamp: ev.Timestamp,
			Delta:     int32(ev.Delta),
			Chars:     model.GetKeyName(uint16(ev.Chars)),
		})
	}

	return &pb.GetCurrentEventsResponse{
		Success: true,
		Events:  pbEvents,
	}, nil
}

// HealthCheck 健康检查
func (s *MacroEngineServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Healthy: true,
		Version: "1.0.0",
	}, nil
}

// convertMacroToProto 转换为 proto
func convertMacroToProto(m *model.Macro) *pb.MacroInfo {
	return &pb.MacroInfo{
		Name:       m.Meta.Name,
		CreatedAt:  m.Meta.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  m.Meta.UpdatedAt.Format("2006-01-02 15:04:05"),
		EventCount: int32(len(m.Events)),
		Script:     m.Script.Code,
	}
}

func main() {
	flag.Parse()

	// 创建服务
	server, err := NewMacroEngineServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterMacroEngineServer(grpcServer, server)

	// 注册健康检查
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// 启用反射（用于调试）
	reflection.Register(grpcServer)

	// 监听端口
	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Macro Engine gRPC server starting on %s", addr)

	// 处理退出信号
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("shutting down server...")
		grpcServer.GracefulStop()
	}()

	// 启动服务器
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
