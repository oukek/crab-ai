package controller

import (
	"oukek/crab-ai/app/service"

	"github.com/gin-gonic/gin"
)

type ChatController struct{}

func (c *ChatController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/chat", c.Chat)
}

// 聊天接口
func (c *ChatController) Chat(ctx *gin.Context) {
	service.Chat(ctx)
}
