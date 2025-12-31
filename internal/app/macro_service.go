package app

import (
	"context"

	"github.com/issueye/macro_operation/internal/model"
	"github.com/issueye/macro_operation/internal/repository"
	"github.com/issueye/macro_operation/internal/service"
)

// MacroService Wails v3 宏服务
type MacroService struct {
	service *service.MacroService
}

// MacroListItem 宏列表项
type MacroListItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	CreatedAt   string `json:"created_at"`
	EventCount  int    `json:"event_count"`
}

// MacroDetail 宏详情
type MacroDetail struct {
	MacroListItem
	ScriptCode string `json:"script_code"`
	Config     Config `json:"config"`
}

// Config 配置
type Config struct {
	PlaybackSpeed float64 `json:"playback_speed"`
	LoopCount     int     `json:"loop_count"`
	AutoSave      bool    `json:"auto_save"`
}

// NewMacroService 创建宏服务
func NewMacroService(repo repository.MacroRepository) *MacroService {
	return &MacroService{
		service: service.NewMacroService(repo),
	}
}

// List 列出所有宏
func (s *MacroService) List(ctx context.Context) ([]*MacroListItem, error) {
	macros, err := s.service.List()
	if err != nil {
		return nil, err
	}

	items := make([]*MacroListItem, len(macros))
	for i, macro := range macros {
		items[i] = &MacroListItem{
			Name:        macro.Meta.Name,
			Description: macro.Meta.Description,
			Version:     macro.Meta.Version,
			CreatedAt:   macro.Meta.CreatedAt.Format("2006-01-02 15:04:05"),
			EventCount:  len(macro.Operations),
		}
	}

	return items, nil
}

// Load 加载宏详情
func (s *MacroService) Load(ctx context.Context, name string) (*MacroDetail, error) {
	macro, err := s.service.Load(name)
	if err != nil {
		return nil, err
	}

	return &MacroDetail{
		MacroListItem: MacroListItem{
			Name:        macro.Meta.Name,
			Description: macro.Meta.Description,
			Version:     macro.Meta.Version,
			CreatedAt:   macro.Meta.CreatedAt.Format("2006-01-02 15:04:05"),
			EventCount:  len(macro.Operations),
		},
		ScriptCode: macro.Script.Code,
		Config: Config{
			PlaybackSpeed: macro.Config.PlaybackSpeed,
			LoopCount:     macro.Config.LoopCount,
			AutoSave:      macro.Config.AutoSave,
		},
	}, nil
}

// Save 保存宏
func (s *MacroService) Save(ctx context.Context, name, description, script string) error {
	macro := model.NewMacro(name)
	macro.Meta.Description = description
	macro.Script.Code = script

	return s.service.Save(macro)
}

// SaveFromRecording 从录制保存宏
func (s *MacroService) SaveFromRecording(ctx context.Context, name, description, script string) error {
	return s.service.SaveFromScript(name, script)
}

// Delete 删除宏
func (s *MacroService) Delete(ctx context.Context, name string) error {
	return s.service.Delete(name)
}

// Exists 检查宏是否存在
func (s *MacroService) Exists(ctx context.Context, name string) (bool, error) {
	return s.service.Exists(name), nil
}

// Rename 重命名宏
func (s *MacroService) Rename(ctx context.Context, oldName, newName string) error {
	return s.service.Rename(oldName, newName)
}
