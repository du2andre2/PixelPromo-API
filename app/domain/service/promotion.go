package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"pixelPromo/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"strings"
	"time"
)

func NewPromotionService(
	rp port.Repository,
	st port.Storage,
	cfg *config.Config,
	log config.Logger,
) port.PromotionHandler {
	return &promotionService{
		rp:  rp,
		st:  st,
		cfg: cfg,
		log: log,
	}
}

type promotionService struct {
	rp  port.Repository
	st  port.Storage
	cfg *config.Config
	log config.Logger
}

func (r *promotionService) CreatePromotion(ctx context.Context, promotion *model.Promotion) error {
	err := r.validPromotion(ctx, promotion)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	promotion.CreatedAt = time.Now()
	promotion.ID = fmt.Sprintf("%d", promotion.CreatedAt.UnixNano())

	if err = r.rp.CreateOrUpdatePromotion(ctx, promotion); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("user created")
	return nil
}

func (r *promotionService) UpdatePromotionImage(ctx context.Context, id string, image io.Reader) error {

	promotion, err := r.rp.GetPromotionByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	if promotion == nil {
		err = errors.New("promotion not found")
		r.log.Error(err.Error())
		return err
	}

	url, err := r.st.UploadPromotionImage(ctx, fmt.Sprintf("%s.jpg", id), image)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	promotion.ImageUrl = url

	if err = r.rp.CreateOrUpdatePromotion(ctx, promotion); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("picture uploaded")
	return nil
}

func (r *promotionService) GetPromotionByID(ctx context.Context, id string) (*model.Promotion, error) {
	promotion, err := r.rp.GetPromotionByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return promotion, nil
}

func (r *promotionService) GetPromotions(ctx context.Context, params *model.PromotionQuery) ([]model.Promotion, error) {
	promotion, err := r.rp.GetPromotionsWithParams(ctx, params)
	if err != nil {
		r.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	return promotion, nil
}

func (r *promotionService) GetPromotionsByCategory(ctx context.Context, category string) ([]model.Promotion, error) {
	promotion, err := r.rp.GetPromotionsByCategory(ctx, category)
	if err != nil {
		r.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	return promotion, nil
}

func (r *promotionService) GetCategories(ctx context.Context) ([]model.Category, error) {
	categories, err := r.rp.GetCategories(ctx)
	if err != nil {
		r.log.Error(err.Error())
		return []model.Category{}, err
	}

	return categories, nil
}

func (r *promotionService) validPromotion(ctx context.Context, promotion *model.Promotion) error {
	if promotion == nil {
		return errors.New("promotion is nil")
	}
	if len(strings.TrimSpace(promotion.Link)) == 0 {
		return errors.New("link is empty")
	}
	if len(strings.TrimSpace(promotion.Description)) == 0 {
		return errors.New("description is empty")
	}
	if len(strings.TrimSpace(promotion.Title)) == 0 {
		return errors.New("title is empty")
	}
	if len(strings.TrimSpace(promotion.UserID)) == 0 {
		return errors.New("userID is empty")
	}

	if len(promotion.Categories) > 0 {
		for _, category := range promotion.Categories {
			if len(strings.TrimSpace(category)) == 0 {
				return errors.New("category name is empty")
			}
		}
	}

	user, err := r.rp.GetUserByID(ctx, promotion.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		err = errors.New("user not found")
		return err
	}

	return nil
}
