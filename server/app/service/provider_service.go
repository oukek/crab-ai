package service

import (
	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service/db"
)

// ProviderService 提供商服务
type ProviderService struct{}

// GetAllProviders 获取所有提供商
func (s *ProviderService) GetAllProviders() ([]model.Provider, error) {
	var providers []model.Provider
	result := db.GetDB().Find(&providers)
	return providers, result.Error
}

// GetProviderByID 根据ID获取提供商
func (s *ProviderService) GetProviderByID(id uint) (model.Provider, error) {
	var provider model.Provider
	result := db.GetDB().First(&provider, id)
	return provider, result.Error
}

// CreateProvider 创建提供商
func (s *ProviderService) CreateProvider(provider *model.Provider) error {
	return db.GetDB().Create(provider).Error
}

// UpdateProvider 更新提供商
func (s *ProviderService) UpdateProvider(provider *model.Provider) error {
	return db.GetDB().Save(provider).Error
}

// DeleteProvider 删除提供商
func (s *ProviderService) DeleteProvider(id uint) error {
	return db.GetDB().Delete(&model.Provider{}, id).Error
}

// GetActiveProviders 获取所有活跃的提供商
func (s *ProviderService) GetActiveProviders() ([]model.Provider, error) {
	var providers []model.Provider
	result := db.GetDB().Where("status = ?", 1).Find(&providers)
	return providers, result.Error
} 