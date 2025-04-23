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

// GetModelsByProviderID 获取指定服务商的所有模型
func (s *ProviderService) GetModelsByProviderID(providerID uint) ([]model.AIModel, error) {
	var models []model.AIModel
	result := db.GetDB().Where("provider_id = ?", providerID).Find(&models)
	return models, result.Error
}

type ProviderDetail struct {
	model.Provider
	Models []model.AIModel `json:"models"`
}

func (s *ProviderService) GetProviderDetailByID(id uint) (ProviderDetail, error) {
	var provider model.Provider
	if err := db.GetDB().First(&provider, id).Error; err != nil {
		return ProviderDetail{}, err
	}
	models, err := s.GetModelsByProviderID(id)
	if err != nil {
		return ProviderDetail{}, err
	}
	return ProviderDetail{Provider: provider, Models: models}, nil
}

// CreateModel 创建模型
func (s *ProviderService) CreateModel(model *model.AIModel) error {
	return db.GetDB().Create(model).Error
}

// UpdateModel 更新模型
func (s *ProviderService) UpdateModel(model *model.AIModel) error {
	return db.GetDB().Save(model).Error
}

// DeleteModel 删除模型
func (s *ProviderService) DeleteModel(id uint) error {
	return db.GetDB().Delete(&model.AIModel{}, id).Error
}

// GetModelByID 根据ID获取模型
func (s *ProviderService) GetModelByID(id uint) (model.AIModel, error) {
	var m model.AIModel
	result := db.GetDB().First(&m, id)
	return m, result.Error
}

// GetAllModels 获取所有模型
func (s *ProviderService) GetAllModels() ([]model.AIModel, error) {
	var models []model.AIModel
	result := db.GetDB().Find(&models)
	return models, result.Error
}
