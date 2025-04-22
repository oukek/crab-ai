package model

import (
	"time"

	"gorm.io/gorm"
)

// AIModelType 定义模型类型常量
const (
	AIModelTypeVision  = "vision"  // 视觉
	AIModelTypeNetwork = "network" // 联网
	AIModelTypeEmbed   = "embed"   // 嵌入
	AIModelTypeReason  = "reason"  // 推理
	AIModelTypeTool    = "tool"    // 工具
)

// AIModel 表示AI模型
type AIModel struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Identifier string         `gorm:"type:varchar(100);not null;unique;comment:模型标识" json:"identifier"`
	Name       string         `gorm:"type:varchar(100);not null;comment:模型名称" json:"name"`
	Group      string         `gorm:"type:varchar(50);comment:模型分组" json:"group"`
	Types      string         `gorm:"type:varchar(255);not null;comment:模型类型,多选用逗号分隔" json:"types"`
	ProviderID uint           `gorm:"not null;comment:服务商编号" json:"provider_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AIModel) TableName() string {
	return "ai_models"
}
