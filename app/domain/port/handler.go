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
	GetInteractionStatisticsByUserID(context.Context, string) (map[string]int, error)
	GetInteractionStatisticsByUserIDWithPromotionID(context.Context, string, string) (map[string]bool, error)

	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	UpdateUser(context.Context, *model.User) error
	DeleteUser(context.Context, string) error
	GetUserByID(context.Context, string) (*model.User, error)
	GetUserRank(context.Context, int) ([]model.User, error)
	Login(context.Context, *model.Login) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	DeletePromotion(context.Context, string) error
	UpdatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetFavoritesPromotionsByUserID(context.Context, string) ([]model.Promotion, error)
	GetPromotions(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
