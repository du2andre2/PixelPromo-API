package route

import (
	"github.com/gin-gonic/gin"
	"pixelPromo/port/controller"
)

type Route interface {
	Setup(severMux *gin.Engine)
}

type route struct {
	controller controller.Controller
}

func NewRoute(
	controller controller.Controller,
) Route {
	return &route{
		controller: controller,
	}
}

func (r *route) Setup(router *gin.Engine) {

	userGroup := router.Group("/users")
	{
		userGroup.POST("/", r.controller.CreateUser)
		userGroup.POST("/picture/:id", r.controller.UpdateUserPicture)
		userGroup.GET("/:id", r.controller.GetUserByID)
	}

	promotionGroup := router.Group("/promotions")
	{
		promotionGroup.POST("/", r.controller.CreatePromotion)
		promotionGroup.POST("/image/:id", r.controller.UpdatePromotionImage)
		promotionGroup.GET("/", r.controller.GetPromotions)
		promotionGroup.GET("/:id", r.controller.GetPromotionByID)
	}

	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("/", r.controller.GetCategories)
	}

	interactionGroup := router.Group("/interactions")
	{
		interactionGroup.POST("/", r.controller.CreateInteraction)
		interactionGroup.GET("/:id", r.controller.GetInteractionByID)
	}

	return
}
