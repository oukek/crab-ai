package model

import (
	"time"

	"gorm.io/gorm"
)

// Prompt 表示大模型常用的指令
// 包含常用字段

type Prompt struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;comment:指令名称" json:"name"`
	Type        string         `gorm:"type:varchar(50);not null;default:'system';comment:指令类型，如system、user等" json:"type"`
	Content     string         `gorm:"type:text;not null;comment:指令内容" json:"content"`
	Description string         `gorm:"type:varchar(255);comment:描述" json:"description"`
	Tags        string         `gorm:"type:varchar(255);comment:标签,逗号分隔" json:"tags"`
	Config      string         `gorm:"type:text;comment:参数配置(JSON)" json:"config"`
	Status      int            `gorm:"type:tinyint;default:1;comment:状态 1:启用 0:禁用" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Prompt) TableName() string {
	return "prompts"
}
