package repository

import (
	"context"
	"io"
	"pixelPromo/domain/model"
)

type InteractionRepository interface {
	CreateInteraction(context.Context, model.PromotionInteraction) error
	GetInteractionByID(context.Context, string) (model.PromotionInteraction, error)
}
type UserRepository interface {
	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(context.Context, string) (*model.User, error)
}
type PromotionRepository interface {
	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetPromotions(context.Context, model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
