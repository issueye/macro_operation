package service

import (
	"fmt"
	"sync"

	"github.com/issueye/macro_operation/internal/model"
	"github.com/issueye/macro_operation/internal/repository"
)

// MacroService 宏管理服务
type MacroService struct {
	repo       repository.MacroRepository
	cacheMutex sync.RWMutex
}

// NewMacroService 创建宏管理服务
func NewMacroService(repo repository.MacroRepository) *MacroService {
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

// Rename 重命名宏
func (s *MacroService) Rename(oldName, newName string) error {
	// 加载旧宏
	macro, err := s.Load(oldName)
	if err != nil {
		return err
	}

	// 检查新名称是否已存在
	if s.repo.Exists(newName) {
		return fmt.Errorf("macro already exists: %s", newName)
	}

	// 删除旧宏
	if err := s.Delete(oldName); err != nil {
		return err
	}

	// 保存新宏
	macro.Meta.Name = newName
	return s.Save(macro)
}

// Exists 检查宏是否存在
func (s *MacroService) Exists(name string) bool {
	return s.repo.Exists(name)
}
