package llm

import (
	"oukek/crab-ai/app/model"
)

// MakeRequest 分发入口
func MakeRequest(chains *ChatChains) error {
	switch chains.Provider.Type {
	case model.ProviderTypeGoogleAIStudio:
		return makeRequestGoogleAIStudio(chains)
	default:
		chains.Set("reqBody", chains.Req)
		return nil
	}
}
