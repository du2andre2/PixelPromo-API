package port

import (
	"context"
	"io"
	"pixelPromo/domain/model"
)

type Handler interface {
	CreateInteraction(context.Context, *model.PromotionInteraction) error
	GetCommentsByPromotionId(context.Context, string) ([]model.PromotionInteraction, error)
	GetInteractionStatisticsByPromotionId(context.Context, string) (map[string]int, error)
	GetInteractionStatisticsByUserId(context.Context, string) (map[string]int, error)
	GetInteractionStatisticsByUserIdWithPromotionId(context.Context, string, string) (map[string]bool, error)

	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	UpdateUser(context.Context, *model.User) error
	DeleteUser(context.Context, string) error
	GetUserById(context.Context, string) (*model.User, error)
	GetUserRank(context.Context, int) ([]model.User, error)
	Login(context.Context, *model.Login) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	DeletePromotion(context.Context, string) error
	UpdatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionById(context.Context, string) (*model.Promotion, error)
	GetFavoritesPromotionsByUserId(context.Context, string) ([]model.Promotion, error)
	GetPromotions(context.Context, *model.PromotionQuery) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}
