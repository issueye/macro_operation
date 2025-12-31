package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/issueye/macro_operation/internal/model"
)

// MacroRepository 宏存储接口
type MacroRepository interface {
	Save(macro *model.Macro) error
	Load(name string) (*model.Macro, error)
	Delete(name string) error
	ListAll() ([]*model.Macro, error)
	Exists(name string) bool
}

// FileMacroRepository 文件存储实现
type FileMacroRepository struct {
	macrosDir  string
	cache      map[string]*model.Macro
	cacheMutex sync.RWMutex
}

// NewFileMacroRepository 创建文件存储
func NewFileMacroRepository(macrosDir string) (*FileMacroRepository, error) {
	// 确保目录存在
	if err := os.MkdirAll(macrosDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create macros directory: %w", err)
	}

	repo := &FileMacroRepository{
		macrosDir: macrosDir,
		cache:     make(map[string]*model.Macro),
	}

	// 加载现有宏
	repo.loadCache()

	return repo, nil
}

// loadCache 加载宏缓存
func (r *FileMacroRepository) loadCache() error {
	files, err := os.ReadDir(r.macrosDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		name := file.Name()[:len(file.Name())-5] // 去掉.json后缀
		macro, err := r.Load(name)
		if err != nil {
			continue
		}

		r.cache[name] = macro
	}

	return nil
}

// Save 保存宏
func (r *FileMacroRepository) Save(macro *model.Macro) error {
	filePath := r.getMacroPath(macro.Meta.Name)

	data, err := json.MarshalIndent(macro, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal macro: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write macro file: %w", err)
	}

	// 更新缓存
	r.cacheMutex.Lock()
	r.cache[macro.Meta.Name] = macro
	r.cacheMutex.Unlock()

	return nil
}

// Load 加载宏
func (r *FileMacroRepository) Load(name string) (*model.Macro, error) {
	filePath := r.getMacroPath(name)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read macro file: %w", err)
	}

	var macro model.Macro
	if err := json.Unmarshal(data, &macro); err != nil {
		return nil, fmt.Errorf("failed to unmarshal macro: %w", err)
	}

	return &macro, nil
}

// Delete 删除宏
func (r *FileMacroRepository) Delete(name string) error {
	filePath := r.getMacroPath(name)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete macro file: %w", err)
	}

	// 从缓存删除
	r.cacheMutex.Lock()
	delete(r.cache, name)
	r.cacheMutex.Unlock()

	return nil
}

// ListAll 列出所有宏
func (r *FileMacroRepository) ListAll() ([]*model.Macro, error) {
	r.cacheMutex.RLock()
	defer r.cacheMutex.RUnlock()

	macros := make([]*model.Macro, 0, len(r.cache))
	for _, macro := range r.cache {
		macros = append(macros, macro)
	}

	return macros, nil
}

// Exists 检查宏是否存在
func (r *FileMacroRepository) Exists(name string) bool {
	r.cacheMutex.RLock()
	defer r.cacheMutex.RUnlock()
	_, exists := r.cache[name]
	return exists
}

// getMacroPath 获取宏文件路径
func (r *FileMacroRepository) getMacroPath(name string) string {
	return filepath.Join(r.macrosDir, name+".json")
}
