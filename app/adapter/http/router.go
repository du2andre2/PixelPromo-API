package http

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Run()
}

type router struct {
	controller *Controller
}

func NewRouter(
	controller *Controller,
) Router {
	return &router{
		controller: controller,
	}
}

func (r *router) Run() {
	gin := gin.Default()
	r.setup(gin)
	gin.Run("localhost:5000")
}

func (r *router) setup(gin *gin.Engine) {

	userGroup := gin.Group("/users")
	{
		userGroup.POST("/", r.controller.CreateUser)
		userGroup.POST("/login", r.controller.Login)
		userGroup.POST("/picture/:id", r.controller.UpdateUserPicture)
		userGroup.GET("/:id", r.controller.GetUserByID)
		userGroup.GET("/rank", r.controller.GetUserRank)
	}

	promotionGroup := gin.Group("/promotions")
	{
		promotionGroup.POST("/", r.controller.CreatePromotion)
		promotionGroup.POST("/image/:id", r.controller.UpdatePromotionImage)
		promotionGroup.GET("/", r.controller.GetPromotions)
		promotionGroup.GET("/:id", r.controller.GetPromotionByID)
		promotionGroup.GET("/favorites/:id", r.controller.GetFavoritesPromotionsByUserID)
	}

	categoryGroup := gin.Group("/categories")
	{
		categoryGroup.GET("/", r.controller.GetCategories)
	}

	interactionGroup := gin.Group("/interactions")
	{
		interactionGroup.POST("/", r.controller.CreateInteraction)
		interactionGroup.GET("/comments/:id", r.controller.GetCommentsByPromotionID)
		interactionGroup.GET("/statistics/:id", r.controller.GetInteractionStatisticsByPromotionID)

	}

	return
}
