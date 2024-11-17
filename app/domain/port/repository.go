package port

import (
	"context"
	"pixelPromo/domain/model"
	"time"
)

type Repository interface {
	CreateOrUpdateInteraction(context.Context, *model.PromotionInteraction) error
	GetInteractionById(context.Context, string) (*model.PromotionInteraction, error)
	DeleteInteraction(context.Context, string) error
	GetInteractionsByPromotionId(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionsByUserId(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionsByUserIdWithPromotionId(context.Context, string, string) ([]model.PromotionInteraction, error)
	GetInteractionsByTypeWithPromotionId(context.Context, model.InteractionType, string) ([]model.PromotionInteraction, error)
	GetInteractionsByTypeWithUserId(context.Context, model.InteractionType, string) ([]model.PromotionInteraction, error)
	CreateOrUpdateUser(context.Context, *model.User) error
	CreateOrUpdateUserScore(context.Context, *model.UserScore) error
	DeleteUser(context.Context, string) error
	GetAllUserScoreByTimeWithUserId(context.Context, string, time.Time) ([]model.UserScore, error)
	GetAllUserScoreByTime(context.Context, time.Time) ([]model.UserScore, error)
	GetUserById(context.Context, string) (*model.User, error)
	GetAllUsers(context.Context) ([]model.User, error)
	GetUserByEmailAndPassword(context.Context, string, string) (*model.User, error)
	CreateOrUpdatePromotion(context.Context, *model.Promotion) error
	DeletePromotion(context.Context, string) error
	GetPromotionById(context.Context, string) (*model.Promotion, error)
	GetPromotionsWithParams(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionsByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
