package configs

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Record   RecordConfig   `mapstructure:"record"`
	Playback PlaybackConfig `mapstructure:"playback"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
	LogPath string `mapstructure:"log_path"`
}

// RecordConfig 录制配置
type RecordConfig struct {
	MaxDuration      int  `mapstructure:"max_duration"`       // 最大录制时长(秒)
	EnableScreenshot bool `mapstructure:"enable_screenshot"`  // 启用截图
	FilterMouseMove  bool `mapstructure:"filter_mouse_move"`  // 过滤鼠标移动事件
	AutoSave         bool `mapstructure:"auto_save"`         // 自动保存
}

// PlaybackConfig 回放配置
type PlaybackConfig struct {
	DefaultSpeed float64 `mapstructure:"default_speed"` // 默认播放速度
	EnableSound  bool    `mapstructure:"enable_sound"`  // 启用声音提示
	StopOnError  bool    `mapstructure:"stop_on_error"` // 错误时停止
}

// StorageConfig 存储配置
type StorageConfig struct {
	MacrosDir      string `mapstructure:"macros_dir"`       // 宏文件目录
	BackupEnabled  bool   `mapstructure:"backup_enabled"`   // 启用备份
	BackupDir      string `mapstructure:"backup_dir"`       // 备份目录
	MaxBackupCount int    `mapstructure:"max_backup_count"` // 最大备份数
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 读取环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 设置默认值
	setDefaults(&config)

	// 创建必要的目录
	createDirectories(&config)

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults(config *Config) {
	if config.App.Name == "" {
		config.App.Name = "Macro Recorder"
	}
	if config.App.Version == "" {
		config.App.Version = "1.0.0"
	}
	if config.Record.MaxDuration == 0 {
		config.Record.MaxDuration = 3600 // 1小时
	}
	if config.Playback.DefaultSpeed == 0 {
		config.Playback.DefaultSpeed = 1.0
	}
	if config.Storage.MacrosDir == "" {
		config.Storage.MacrosDir = "./macros"
	}
	if config.Storage.BackupDir == "" {
		config.Storage.BackupDir = "./backups"
	}
	if config.App.LogPath == "" {
		config.App.LogPath = "./logs"
	}
}

// createDirectories 创建必要的目录
func createDirectories(config *Config) {
	dirs := []string{
		config.Storage.MacrosDir,
		config.App.LogPath,
		config.Storage.BackupDir,
	}

	for _, dir := range dirs {
		if dir != "" {
			os.MkdirAll(dir, 0755)
		}
	}
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	// 检查当前目录
	if _, err := os.Stat("configs/config.yaml"); err == nil {
		return "configs/config.yaml"
	}

	// 检查可执行文件目录
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		configPath := filepath.Join(exeDir, "configs", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// 默认配置路径
	return "configs/config.yaml"
}

// Save 保存配置
func Save(configPath string, config *Config) error {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 设置配置值
	if err := v.MergeConfigMap(configToMap(config)); err != nil {
		return err
	}

	// 写入配置文件
	return v.WriteConfigAs(configPath)
}

// configToMap 将配置转换为 map
func configToMap(config *Config) map[string]interface{} {
	return map[string]interface{}{
		"app": map[string]interface{}{
			"name":     config.App.Name,
			"version":  config.App.Version,
			"debug":    config.App.Debug,
			"log_path": config.App.LogPath,
		},
		"record": map[string]interface{}{
			"max_duration":       config.Record.MaxDuration,
			"enable_screenshot":  config.Record.EnableScreenshot,
			"filter_mouse_move":  config.Record.FilterMouseMove,
			"auto_save":          config.Record.AutoSave,
		},
		"playback": map[string]interface{}{
			"default_speed": config.Playback.DefaultSpeed,
			"enable_sound":  config.Playback.EnableSound,
			"stop_on_error": config.Playback.StopOnError,
		},
		"storage": map[string]interface{}{
			"macros_dir":       config.Storage.MacrosDir,
			"backup_enabled":   config.Storage.BackupEnabled,
			"backup_dir":       config.Storage.BackupDir,
			"max_backup_count": config.Storage.MaxBackupCount,
		},
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "Macro Recorder",
			Version: "1.0.0",
			Debug:   false,
			LogPath: "./logs",
		},
		Record: RecordConfig{
			MaxDuration:      3600,
			EnableScreenshot: false,
			FilterMouseMove:  true,
			AutoSave:         false,
		},
		Playback: PlaybackConfig{
			DefaultSpeed: 1.0,
			EnableSound:  false,
			StopOnError:  true,
		},
		Storage: StorageConfig{
			MacrosDir:      "./macros",
			BackupEnabled:  true,
			BackupDir:      "./backups",
			MaxBackupCount: 10,
		},
	}
}
