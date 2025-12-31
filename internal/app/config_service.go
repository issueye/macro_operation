package app

import (
	"context"

	"github.com/issueye/macro_operation/configs"
)

// ConfigService Wails v3 配置服务
type ConfigService struct {
	config *configs.Config
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Debug   bool   `json:"debug"`
}

// RecordConfig 录制配置
type RecordConfig struct {
	MaxDuration      int    `json:"max_duration"`
	EnableScreenshot bool   `json:"enable_screenshot"`
	FilterMouseMove  bool   `json:"filter_mouse_move"`
	RecordHotkey     string `json:"record_hotkey"`
	StopHotkey       string `json:"stop_hotkey"`
	AutoSave         bool   `json:"auto_save"`
}

// PlaybackConfig 播放配置
type PlaybackConfig struct {
	DefaultSpeed  float64 `json:"default_speed"`
	EnableSound   bool    `json:"enable_sound"`
	StopOnError   bool    `json:"stop_on_error"`
	ConfirmBefore bool    `json:"confirm_before"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	MacrosDir      string `json:"macros_dir"`
	BackupEnabled  bool   `json:"backup_enabled"`
	BackupDir      string `json:"backup_dir"`
	MaxBackupCount int    `json:"max_backup_count"`
}

// NewConfigService 创建配置服务
func NewConfigService(config *configs.Config) *ConfigService {
	return &ConfigService{config: config}
}

// GetAppConfig 获取应用配置
func (s *ConfigService) GetAppConfig(ctx context.Context) (*AppConfig, error) {
	return &AppConfig{
		Name:    s.config.App.Name,
		Version: s.config.App.Version,
		Debug:   s.config.App.Debug,
	}, nil
}

// GetRecordConfig 获取录制配置
func (s *ConfigService) GetRecordConfig(ctx context.Context) (*RecordConfig, error) {
	return &RecordConfig{
		MaxDuration:      s.config.Record.MaxDuration,
		EnableScreenshot: s.config.Record.EnableScreenshot,
		FilterMouseMove:  s.config.Record.FilterMouseMove,
		AutoSave:         s.config.Record.AutoSave,
	}, nil
}

// UpdateRecordConfig 更新录制配置
func (s *ConfigService) UpdateRecordConfig(ctx context.Context, config *RecordConfig) error {
	s.config.Record.MaxDuration = config.MaxDuration
	s.config.Record.EnableScreenshot = config.EnableScreenshot
	s.config.Record.FilterMouseMove = config.FilterMouseMove
	s.config.Record.AutoSave = config.AutoSave

	// 保存到文件
	return configs.Save(configs.GetConfigPath(), s.config)
}

// GetPlaybackConfig 获取播放配置
func (s *ConfigService) GetPlaybackConfig(ctx context.Context) (*PlaybackConfig, error) {
	return &PlaybackConfig{
		DefaultSpeed: s.config.Playback.DefaultSpeed,
		EnableSound:  s.config.Playback.EnableSound,
		StopOnError:  s.config.Playback.StopOnError,
	}, nil
}

// UpdatePlaybackConfig 更新播放配置
func (s *ConfigService) UpdatePlaybackConfig(ctx context.Context, config *PlaybackConfig) error {
	s.config.Playback.DefaultSpeed = config.DefaultSpeed
	s.config.Playback.EnableSound = config.EnableSound
	s.config.Playback.StopOnError = config.StopOnError

	return configs.Save(configs.GetConfigPath(), s.config)
}

// GetStorageConfig 获取存储配置
func (s *ConfigService) GetStorageConfig(ctx context.Context) (*StorageConfig, error) {
	return &StorageConfig{
		MacrosDir:      s.config.Storage.MacrosDir,
		BackupEnabled:  s.config.Storage.BackupEnabled,
		BackupDir:      s.config.Storage.BackupDir,
		MaxBackupCount: s.config.Storage.MaxBackupCount,
	}, nil
}
