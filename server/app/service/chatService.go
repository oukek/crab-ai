package service

import (
	llm "oukek/crab-ai/app/service/llm"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/gin-gonic/gin"
)

// ChatRequest 聊天请求结构体
type ChatRequest struct {
	Content string `json:"content"`  // 用户输入内容
	ModelID uint   `json:"model_id"` // 选择的模型id
}

// Chat 聊天主方法
func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误"})
		return
	}

	providerService := &ProviderService{}
	// 查询模型
	model, err := providerService.GetModelByID(req.ModelID)
	if err != nil {
		c.JSON(400, gin.H{"error": "模型不存在"})
		return
	}

	// 查询服务商
	provider, err := providerService.GetProviderByID(model.ProviderID)
	if err != nil {
		c.JSON(400, gin.H{"error": "服务商不存在"})
		return
	}

	// 初始化 llmReq
	llmReq := &_type.Request{
		Messages: []_type.Message{
			{
				Role:    _type.UserRole,
				Content: req.Content,
			},
		},
		Model: string(model.Identifier),
	}

	resp, err := llm.SimpleChatWithReq(c, llmReq, provider)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}
