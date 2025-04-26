package llm

import (
	"fmt"
	"sync"

	"oukek/crab-ai/app/model"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[model.ProviderType]LLMProvider)
)

// RegisterProvider 注册一个 LLMProvider 实现
// 如果需要支持新的 Provider 类型，需要在此处添加 case
// 为了避免循环依赖，具体的 Provider 实现应该在它们自己的包（如果需要）
// 或者在 llm 包内部定义，并在 init 函数中调用 RegisterProvider
func RegisterProvider(providerType model.ProviderType, provider LLMProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("RegisterProvider: provider is nil")
	}
	if _, dup := providers[providerType]; dup {
		panic("RegisterProvider: Register called twice for provider type " + string(providerType))
	}
	providers[providerType] = provider
}

// init 在包加载时自动注册已知的 Provider
func init() {
	// 注册 Google AI Studio Provider
	RegisterProvider(model.ProviderTypeGoogleAIStudio, NewGoogleAIStudioProvider())

	// TODO: 在此处添加其他 Provider 的注册
	// 例如: RegisterProvider(model.ProviderTypeOpenAI, NewOpenAIProvider())
}

// GetProvider 根据 Provider 类型获取对应的 LLMProvider 实例
func GetProvider(providerType model.ProviderType) (LLMProvider, error) {
	providersMu.RLock()
	provider, ok := providers[providerType]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("不支持的 Provider 类型: %s", providerType)
	}
	return provider, nil
}
