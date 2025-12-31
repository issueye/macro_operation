package model

import "time"

// Macro 宏模型
type Macro struct {
	Meta       MetaInfo                   `json:"meta"`
	Config     Config                     `json:"config"`
	Script     Script                     `json:"script"`
	Operations []map[string]interface{}   `json:"operations,omitempty"`
}

// MetaInfo 元信息
type MetaInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      string    `json:"author,omitempty"`
}

// Config 配置信息
type Config struct {
	PlaybackSpeed float64 `json:"playback_speed"`
	LoopCount     int     `json:"loop_count"`
	AutoSave      bool    `json:"auto_save"`
}

// Script 脚本信息
type Script struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// NewMacro 创建新宏
func NewMacro(name string) *Macro {
	now := time.Now()
	return &Macro{
		Meta: MetaInfo{
			Name:      name,
			Version:   "1.0.0",
			CreatedAt: now,
			UpdatedAt: now,
		},
		Config: Config{
			PlaybackSpeed: 1.0,
			LoopCount:     1,
			AutoSave:      true,
		},
		Script: Script{
			Language: "javascript",
		},
	}
}
