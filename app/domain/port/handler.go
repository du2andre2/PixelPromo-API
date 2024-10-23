package port

import (
	"context"
	"io"
	"pixelPromo/domain/model"
)

type Handler interface {
	CreateInteraction(context.Context, *model.PromotionInteraction) error
	GetCommentsByPromotionID(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionStatisticsByPromotionID(context.Context, string) (map[string]int, error)

	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(context.Context, string) (*model.User, error)
	GetUserRank(context.Context, int) ([]model.User, error)
	Login(context.Context, *model.Login) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetFavoritesPromotionsByUserID(context.Context, string) ([]model.Promotion, error)
	GetPromotions(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
