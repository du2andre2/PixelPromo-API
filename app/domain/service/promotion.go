package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"pixelPromo/domain/model"
	"strings"
	"time"
)

func (s *service) CreatePromotion(ctx context.Context, promotion *model.Promotion) error {
	err := s.validPromotion(ctx, promotion)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	promotion.CreatedAt = time.Now()
	promotion.ID = fmt.Sprintf("%d", promotion.CreatedAt.UnixNano())

	if err = s.rp.CreateOrUpdatePromotion(ctx, promotion); err != nil {
		s.log.Error(err.Error())
		return err
	}

	interaction := model.PromotionInteraction{
		ID:              promotion.ID,
		PromotionID:     promotion.ID,
		OwnerUserID:     promotion.UserID,
		UserID:          promotion.UserID,
		InteractionType: model.Create,
		CreatedAt:       promotion.CreatedAt,
	}

	if err = s.CreateInteraction(ctx, &interaction); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("user created")
	return nil
}

func (s *service) UpdatePromotionImage(ctx context.Context, id string, image io.Reader) error {

	promotion, err := s.rp.GetPromotionByID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	if promotion == nil {
		err = errors.New("promotion not found")
		s.log.Error(err.Error())
		return err
	}

	url, err := s.st.UploadPromotionImage(ctx, fmt.Sprintf("%s.jpg", id), image)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	promotion.ImageUrl = url

	if err = s.rp.CreateOrUpdatePromotion(ctx, promotion); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("picture uploaded")
	return nil
}

func (s *service) GetPromotionByID(ctx context.Context, id string) (*model.Promotion, error) {
	promotion, err := s.rp.GetPromotionByID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return promotion, nil
}

func (s *service) GetPromotions(ctx context.Context, params *model.PromotionQuery) ([]model.Promotion, error) {
	promotion, err := s.rp.GetPromotionsWithParams(ctx, params)
	if err != nil {
		s.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	return promotion, nil
}

func (s *service) GetFavoritesPromotionsByUserID(ctx context.Context, userID string) ([]model.Promotion, error) {
	interactions, err := s.rp.GetInteractionsByTypeWithUserID(ctx, model.Favorite, userID)
	if err != nil {
		s.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	promotions := make([]model.Promotion, 0)
	for _, interaction := range interactions {
		promotion, err := s.rp.GetPromotionByID(ctx, interaction.PromotionID)
		if err != nil {
			s.log.Error(err.Error())
			return []model.Promotion{}, err
		}

		promotions = append(promotions, *promotion)

	}

	return promotions, nil

}

func (s *service) GetPromotionsByCategory(ctx context.Context, category string) ([]model.Promotion, error) {
	promotion, err := s.rp.GetPromotionsByCategory(ctx, category)
	if err != nil {
		s.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	return promotion, nil
}

func (s *service) GetCategories(ctx context.Context) ([]model.Category, error) {
	categories, err := s.rp.GetCategories(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return []model.Category{}, err
	}

	return categories, nil
}

func (s *service) validPromotion(ctx context.Context, promotion *model.Promotion) error {
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

	user, err := s.rp.GetUserByID(ctx, promotion.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		err = errors.New("user not found")
		return err
	}

	return nil
}
