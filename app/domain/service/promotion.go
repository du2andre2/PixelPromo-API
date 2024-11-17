package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
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

	promotion.DiscountBadge = math.Round(((promotion.OriginalPrice - promotion.DiscountedPrice) / promotion.OriginalPrice) * 100)

	promotion.CreatedAt = time.Now()
	promotion.Id = fmt.Sprintf("%d", promotion.CreatedAt.UnixNano())

	if err = s.rp.CreateOrUpdatePromotion(ctx, promotion); err != nil {
		s.log.Error(err.Error())
		return err
	}

	interaction := model.PromotionInteraction{
		Id:              promotion.Id,
		PromotionId:     promotion.Id,
		OwnerUserId:     promotion.UserId,
		UserId:          promotion.UserId,
		InteractionType: model.Create,
		CreatedAt:       promotion.CreatedAt,
	}

	if err = s.CreateInteraction(ctx, &interaction); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("promotion created")
	return nil
}

func (s *service) DeletePromotion(ctx context.Context, promotionId string) error {

	if err := s.rp.DeletePromotion(ctx, promotionId); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("promotion deleted")
	return nil
}
func (s *service) UpdatePromotion(ctx context.Context, newPromotion *model.Promotion) error {
	err := s.validPromotion(ctx, newPromotion)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	if newPromotion.Id == "" {
		err = errors.New("promotion id is empty")
		s.log.Error(err.Error())
		return err
	}

	promotion, err := s.rp.GetPromotionById(ctx, newPromotion.Id)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	if promotion == nil {
		err = errors.New("promotion not found")
		s.log.Error(err.Error())
		return err
	}

	if err = s.rp.CreateOrUpdatePromotion(ctx, newPromotion); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("promotion updated")
	return nil
}

func (s *service) UpdatePromotionImage(ctx context.Context, id string, image io.Reader) error {

	promotion, err := s.rp.GetPromotionById(ctx, id)
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

func (s *service) GetPromotionById(ctx context.Context, id string) (*model.Promotion, error) {
	promotion, err := s.rp.GetPromotionById(ctx, id)
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

func (s *service) GetFavoritesPromotionsByUserId(ctx context.Context, userId string) ([]model.Promotion, error) {
	interactions, err := s.rp.GetInteractionsByTypeWithUserId(ctx, model.Favorite, userId)
	if err != nil {
		s.log.Error(err.Error())
		return []model.Promotion{}, err
	}

	promotions := make([]model.Promotion, 0)
	for _, interaction := range interactions {
		promotion, err := s.rp.GetPromotionById(ctx, interaction.PromotionId)
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
	if len(strings.TrimSpace(promotion.Title)) == 0 {
		return errors.New("title is empty")
	}
	if len(strings.TrimSpace(promotion.UserId)) == 0 {
		return errors.New("userId is empty")
	}

	if len(promotion.Categories) > 0 {
		for _, category := range promotion.Categories {
			if len(strings.TrimSpace(category)) == 0 {
				return errors.New("category name is empty")
			}
		}
	}

	user, err := s.rp.GetUserById(ctx, promotion.UserId)
	if err != nil {
		return err
	}

	if user == nil {
		err = errors.New("user not found")
		return err
	}

	return nil
}
