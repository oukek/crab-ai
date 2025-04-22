package controller

import (
	"net/http"
	"strconv"

	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service"

	"github.com/gin-gonic/gin"
)

var promptService = &service.PromptService{}

type PromptController struct{}

func (c *PromptController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/prompts", c.GetAll)
	r.GET("/prompts/:id", c.GetByID)
	r.POST("/prompts", c.Create)
	r.PUT("/prompts/:id", c.Update)
	r.DELETE("/prompts/:id", c.Delete)
}

func (c *PromptController) GetAll(ctx *gin.Context) {
	prompts, err := promptService.GetAllPrompts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prompts)
}

func (c *PromptController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	prompt, err := promptService.GetPromptByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prompt)
}

func (c *PromptController) Create(ctx *gin.Context) {
	var prompt model.Prompt
	if err := ctx.ShouldBindJSON(&prompt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := promptService.CreatePrompt(&prompt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prompt)
}

func (c *PromptController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var prompt model.Prompt
	if err := ctx.ShouldBindJSON(&prompt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prompt.ID = uint(id)
	if err := promptService.UpdatePrompt(&prompt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, prompt)
}

func (c *PromptController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := promptService.DeletePrompt(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
