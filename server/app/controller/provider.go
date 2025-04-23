package controller

import (
	"net/http"
	"strconv"

	"oukek/crab-ai/app/model"
	"oukek/crab-ai/app/service"

	"github.com/gin-gonic/gin"
)

var providerService = &service.ProviderService{}

type ProviderController struct{}

func (c *ProviderController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/providers", c.GetAll)
	r.GET("/providers/:id", c.GetByID)
	r.POST("/providers", c.Create)
	r.PUT("/providers/:id", c.Update)
	r.DELETE("/providers/:id", c.Delete)

	// 模型相关路由
	r.GET("/models", c.GetAllModels)
	r.GET("/models/:id", c.GetModelByID)
	r.POST("/models", c.CreateModel)
	r.PUT("/models/:id", c.UpdateModel)
	r.DELETE("/models/:id", c.DeleteModel)
}

func (c *ProviderController) GetAll(ctx *gin.Context) {
	providers, err := providerService.GetAllProviders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, providers)
}

func (c *ProviderController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	provider, err := providerService.GetProviderDetailByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, provider)
}

func (c *ProviderController) Create(ctx *gin.Context) {
	var provider model.Provider
	if err := ctx.ShouldBindJSON(&provider); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := providerService.CreateProvider(&provider); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, provider)
}

func (c *ProviderController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var provider model.Provider
	if err := ctx.ShouldBindJSON(&provider); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	provider.ID = uint(id)
	if err := providerService.UpdateProvider(&provider); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, provider)
}

func (c *ProviderController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := providerService.DeleteProvider(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetAllModels 获取所有模型
func (c *ProviderController) GetAllModels(ctx *gin.Context) {
	models, err := providerService.GetAllModels()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models)
}

// GetModelByID 获取单个模型
func (c *ProviderController) GetModelByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	model, err := providerService.GetModelByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, model)
}

// CreateModel 新建模型
func (c *ProviderController) CreateModel(ctx *gin.Context) {
	var m model.AIModel
	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := providerService.CreateModel(&m); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, m)
}

// UpdateModel 更新模型
func (c *ProviderController) UpdateModel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var m model.AIModel
	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := providerService.UpdateModel(&m); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, m)
}

// DeleteModel 删除模型
func (c *ProviderController) DeleteModel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := providerService.DeleteModel(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
