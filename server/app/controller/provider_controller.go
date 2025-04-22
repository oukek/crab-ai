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
