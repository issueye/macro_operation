package model

import "time"

// Macro 宏模型
type Macro struct {
	Meta   MacroMeta   `json:"meta"`
	Script MacroScript `json:"script"`
	Events []Event     `json:"events"`
}

// MacroMeta 宏元数据
type MacroMeta struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MacroScript 宏脚本
type MacroScript struct {
	Code   string `json:"code"`
	Lang   string `json:"lang"`
}

// NewMacro 创建新宏
func NewMacro(name string) *Macro {
	now := time.Now()
	return &Macro{
		Meta: MacroMeta{
			Name:      name,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Script: MacroScript{
			Lang: "javascript",
		},
		Events: []Event{},
	}
}
