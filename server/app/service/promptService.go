package service

import (
	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service/db"
)

// PromptService 指令服务
type PromptService struct{}

// GetAllPrompts 获取所有指令
func (s *PromptService) GetAllPrompts() ([]model.Prompt, error) {
	var prompts []model.Prompt
	result := db.GetDB().Find(&prompts)
	return prompts, result.Error
}

// GetPromptByID 根据ID获取指令
func (s *PromptService) GetPromptByID(id uint) (model.Prompt, error) {
	var prompt model.Prompt
	result := db.GetDB().First(&prompt, id)
	return prompt, result.Error
}

// CreatePrompt 创建指令
func (s *PromptService) CreatePrompt(prompt *model.Prompt) error {
	return db.GetDB().Create(prompt).Error
}

// UpdatePrompt 更新指令
func (s *PromptService) UpdatePrompt(prompt *model.Prompt) error {
	return db.GetDB().Save(prompt).Error
}

// DeletePrompt 删除指令
func (s *PromptService) DeletePrompt(id uint) error {
	return db.GetDB().Delete(&model.Prompt{}, id).Error
}

// GetActivePrompts 获取所有启用的指令
func (s *PromptService) GetActivePrompts() ([]model.Prompt, error) {
	var prompts []model.Prompt
	result := db.GetDB().Where("status = ?", 1).Find(&prompts)
	return prompts, result.Error
}
