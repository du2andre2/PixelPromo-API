package port

import (
	"context"
	"pixelPromo/domain/model"
)

type Repository interface {
	CreateOrUpdateInteraction(context.Context, *model.PromotionInteraction) error
	GetInteractionByID(context.Context, string) (*model.PromotionInteraction, error)
	GetInteractionsByPromotionID(context.Context, string) ([]model.PromotionInteraction, error)
	GetCommentsByPromotionID(context.Context, string) ([]model.PromotionInteraction, error)
	CreateOrUpdateUser(context.Context, *model.User) error
	CreateOrUpdateUserScore(context.Context, *model.UserScore) error
	GetAllUserScoreByRange(context.Context, string, int) ([]model.UserScore, error)
	GetUserByID(context.Context, string) (*model.User, error)
	CreateOrUpdatePromotion(context.Context, *model.Promotion) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetPromotionsWithParams(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionsByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
