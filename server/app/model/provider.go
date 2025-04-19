package model

import (
	"time"

	"gorm.io/gorm"
)

// Provider 提供商模型
type Provider struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;comment:提供商名称" json:"name"`
	ApiKey    string         `gorm:"type:varchar(255);comment:API密钥" json:"api_key"`
	ApiSecret string         `gorm:"type:varchar(255);comment:API密钥" json:"api_secret"`
	BaseUrl   string         `gorm:"type:varchar(255);comment:基础URL" json:"base_url"`
	Status    int            `gorm:"type:tinyint;default:1;comment:状态 1:启用 0:禁用" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Provider) TableName() string {
	return "providers"
} 