package llm

import (
	"oukek/crab-ai/app/service/llm/provider"
)

// 根据不同的模型，将其转为json数据
func makeRequest(chains *ChatChains) error {
	return provider.MakeRequest(chains)
}
