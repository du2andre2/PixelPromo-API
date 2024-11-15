package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("my_secret_key")

type Router interface {
	Run()
}

type router struct {
	controller *Controller
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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

	gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173/"}, // Porta do frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	gin.POST("/auth", r.controller.Login)
	gin.POST("/users", r.controller.CreateUser)

	userGroup := gin.Group("/users")
	//userGroup.Use(authMiddleware())
	{
		userGroup.POST("/picture/:id", r.controller.UpdateUserPicture)
		userGroup.PATCH("", r.controller.UpdateUser)
		userGroup.DELETE(":id", r.controller.DeleteUser)
		userGroup.GET(":id", r.controller.GetUserByID)
		userGroup.GET("/rank", r.controller.GetUserRank)
	}

	promotionGroup := gin.Group("/promotions")
	//promotionGroup.Use(authMiddleware())
	{
		promotionGroup.POST("", r.controller.CreatePromotion)
		promotionGroup.DELETE(":id", r.controller.DeletePromotion)
		promotionGroup.PATCH("", r.controller.UpdatePromotion)
		promotionGroup.POST("/image/:id", r.controller.UpdatePromotionImage)
		promotionGroup.GET("", r.controller.GetPromotions) // queryParams: []category, search
		promotionGroup.GET(":id", r.controller.GetPromotionByID)
		promotionGroup.GET("/favorites/:id", r.controller.GetFavoritesPromotionsByUserID)
	}

	categoryGroup := gin.Group("/categories")
	//categoryGroup.Use(authMiddleware())
	{
		categoryGroup.GET("", r.controller.GetCategories)
	}

	interactionGroup := gin.Group("/interactions")
	//interactionGroup.Use(authMiddleware())
	{
		interactionGroup.POST("", r.controller.CreateInteraction)
		interactionGroup.GET("/comments/:id", r.controller.GetCommentsByPromotionID)
		interactionGroup.GET("/statistics/:id", r.controller.GetInteractionStatisticsByPromotionID)
		interactionGroup.GET("/user-statistics/:id", r.controller.GetInteractionStatisticsByUserID)
		interactionGroup.GET("/promotion-user-statistics", r.controller.GetInteractionStatisticsByUserIDWithPromotionID)
	}

	return
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
