package port

import (
	"context"
	"io"
	"pixelPromo/domain/model"
)

type InteractionHandler interface {
	CreateOrUpdateInteraction(context.Context, *model.PromotionInteraction) error
	GetInteractionByID(context.Context, string) (*model.PromotionInteraction, error)
	GetCommentsByPromotionID(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionsCountersByPromotionID(context.Context, string) (map[string]int, error)
}

type UserHandler interface {
	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(context.Context, string) (*model.User, error)
	GetUserRank(context.Context) ([]model.User, error)
	Login(context.Context, *model.Login) (*model.User, error)
}
type PromotionHandler interface {
	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetPromotions(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionsByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
