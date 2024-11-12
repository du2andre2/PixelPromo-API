package port

import (
	"context"
	"pixelPromo/domain/model"
	"time"
)

type Repository interface {
	CreateOrUpdateInteraction(context.Context, *model.PromotionInteraction) error
	GetInteractionByID(context.Context, string) (*model.PromotionInteraction, error)
	DeleteInteraction(context.Context, string) error
	GetInteractionsByPromotionID(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionsByUserID(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionsByUserIDWithPromotionID(context.Context, string, string) ([]model.PromotionInteraction, error)
	GetInteractionsByTypeWithPromotionID(context.Context, model.InteractionType, string) ([]model.PromotionInteraction, error)
	GetInteractionsByTypeWithUserID(context.Context, model.InteractionType, string) ([]model.PromotionInteraction, error)
	CreateOrUpdateUser(context.Context, *model.User) error
	CreateOrUpdateUserScore(context.Context, *model.UserScore) error
	GetAllUserScoreByTimeWithUserId(context.Context, string, time.Time) ([]model.UserScore, error)
	GetAllUserScoreByTime(context.Context, time.Time) ([]model.UserScore, error)
	GetUserByID(context.Context, string) (*model.User, error)
	GetAllUsers(context.Context) ([]model.User, error)
	GetUserByEmailAndPassword(context.Context, string, string) (*model.User, error)
	CreateOrUpdatePromotion(context.Context, *model.Promotion) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetPromotionsWithParams(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionsByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
