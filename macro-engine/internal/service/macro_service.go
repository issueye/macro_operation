package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"macro-engine/internal/model"
)

// MacroRepository 宏仓储接口
type MacroRepository interface {
	Save(macro *model.Macro) error
	Load(name string) (*model.Macro, error)
	Delete(name string) error
	ListAll() ([]*model.Macro, error)
	Exists(name string) bool
}

// FileMacroRepository 文件仓储实现
type FileMacroRepository struct {
	BasePath string
	mutex    sync.RWMutex
}

// NewFileMacroRepository 创建文件仓储
func NewFileMacroRepository(basePath string) (*FileMacroRepository, error) {
	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	return &FileMacroRepository{
		BasePath: basePath,
	}, nil
}

// GetFilePath 获取文件路径
func (r *FileMacroRepository) GetFilePath(name string) string {
	return filepath.Join(r.BasePath, name+".json")
}

// Save 保存宏
func (r *FileMacroRepository) Save(macro *model.Macro) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	macro.Meta.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(macro, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal macro: %w", err)
	}

	filePath := r.GetFilePath(macro.Meta.Name)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load 加载宏
func (r *FileMacroRepository) Load(name string) (*model.Macro, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	filePath := r.GetFilePath(name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var macro model.Macro
	if err := json.Unmarshal(data, &macro); err != nil {
		return nil, fmt.Errorf("failed to unmarshal macro: %w", err)
	}

	return &macro, nil
}

// Delete 删除宏
func (r *FileMacroRepository) Delete(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	filePath := r.GetFilePath(name)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListAll 列出所有宏
func (r *FileMacroRepository) ListAll() ([]*model.Macro, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	entries, err := os.ReadDir(r.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	macros := make([]*model.Macro, 0)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // 去掉 .json 后缀
		macro, err := r.Load(name)
		if err != nil {
			continue // 跳过加载失败的宏
		}
		macros = append(macros, macro)
	}

	return macros, nil
}

// Exists 检查宏是否存在
func (r *FileMacroRepository) Exists(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	filePath := r.GetFilePath(name)
	_, err := os.Stat(filePath)
	return err == nil
}

// MacroService 宏管理服务
type MacroService struct {
	repo       MacroRepository
	cacheMutex sync.RWMutex
}

// NewMacroService 创建宏管理服务
func NewMacroService(repo MacroRepository) *MacroService {
	return &MacroService{
		repo: repo,
	}
}

// Save 保存宏
func (s *MacroService) Save(macro *model.Macro) error {
	if err := s.repo.Save(macro); err != nil {
		return fmt.Errorf("failed to save macro: %w", err)
	}
	return nil
}

// SaveFromScript 从脚本保存宏
func (s *MacroService) SaveFromScript(name, script string) error {
	macro := model.NewMacro(name)
	macro.Script.Code = script
	macro.Events = nil // 脚本模式下不需要事件

	return s.Save(macro)
}

// Load 加载宏
func (s *MacroService) Load(name string) (*model.Macro, error) {
	macro, err := s.repo.Load(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load macro: %w", err)
	}
	return macro, nil
}

// Delete 删除宏
func (s *MacroService) Delete(name string) error {
	if err := s.repo.Delete(name); err != nil {
		return fmt.Errorf("failed to delete macro: %w", err)
	}
	return nil
}

// List 列出所有宏
func (s *MacroService) List() ([]*model.Macro, error) {
	macros, err := s.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list macros: %w", err)
	}
	return macros, nil
}

// Exists 检查宏是否存在
func (s *MacroService) Exists(name string) bool {
	return s.repo.Exists(name)
}
