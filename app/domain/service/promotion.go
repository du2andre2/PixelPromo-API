package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port/handler"
	"strings"
	"time"
)

func NewPromotionService(
	db aws.DynamoDB,
	s3 aws.BucketS3,
	cfg *config.Config,
	log config.Logger,
) handler.PromotionHandler {
	return &promotionService{
		db:  db,
		s3:  s3,
		cfg: cfg,
		log: log,
	}
}

type promotionService struct {
	db  aws.DynamoDB
	s3  aws.BucketS3
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

	if err = r.db.PutItem(ctx,
		&aws.PutItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			BodyItem:  promotion,
		}); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("user created")
	return nil
}

func (r *promotionService) UpdatePromotionImage(ctx context.Context, id string, image io.Reader) error {
	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     id,
					ValueType: aws.String},
			},
		})

	if out == nil || out.Item == nil {
		err = errors.New("promotion not found")
		r.log.Error(err.Error())
		return err
	}

	outS3, err := r.s3.PutFile(ctx,
		&aws.PutFileInput{
			BucketName: r.cfg.Viper.GetString("aws.s3.buckets.promotion-images"),
			FileName:   fmt.Sprintf("%s.jpg", id),
			BodyItem:   image,
		},
	)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	if outS3 == nil {
		err = errors.New("error uploading file")
		r.log.Error(err.Error())
		return err
	}

	user := &model.Promotion{}
	err = json.Unmarshal(out.Item, user)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	user.ImageUrl = outS3.UrlFile

	if err = r.db.PutItem(ctx,
		&aws.PutItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			BodyItem:  user,
		}); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("picture uploaded")
	return nil
}

func (r *promotionService) GetPromotionByID(ctx context.Context, id string) (*model.Promotion, error) {
	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     id,
					ValueType: aws.String},
			},
		})

	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	if out == nil || out.Item == nil {
		return nil, nil
	}

	promotion := &model.Promotion{}
	err = json.Unmarshal(out.Item, promotion)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return promotion, nil
}

func (r *promotionService) GetPromotions(ctx context.Context, params model.PromotionQuery) ([]model.Promotion, error) {
	var conditions []aws.ConditionParam
	for _, category := range params.Categories {
		conditions = append(conditions, aws.ConditionParam{
			Names:         []string{"categories"},
			Value:         category,
			OperationType: aws.Contains,
		})
	}
	if params.Search != "" {
		conditions = append(conditions,
			aws.ConditionParam{
				Names:         []string{"userId", "title", "description", "platform", "categories", "link"},
				Value:         params.Search,
				OperationType: aws.Contains,
			},
		)
	}

	out, err := r.db.ScanItem(ctx,
		&aws.ScanItemInput{
			TableName:          r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			Conditions:         conditions,
			ConditionConstrain: aws.ANYCondition,
		})

	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	if out == nil || len(out.Items) == 0 {
		return nil, nil
	}

	var promotions []model.Promotion

	for _, out := range out.Items {
		promotion := model.Promotion{}
		err = json.Unmarshal(out.Item, &promotion)
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		promotions = append(promotions, promotion)
	}

	return promotions, nil
}

func (r *promotionService) GetPromotionByCategory(ctx context.Context, category string) ([]model.Promotion, error) {
	out, err := r.db.ScanItem(ctx,
		&aws.ScanItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			Conditions: []aws.ConditionParam{
				{
					Names:         []string{"categories"},
					Value:         category,
					OperationType: aws.Contains,
				},
			},
		})

	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	if out == nil || out.Items == nil || len(out.Items) == 0 {
		return nil, nil
	}

	var promotions []model.Promotion

	for _, out := range out.Items {
		promotion := model.Promotion{}
		err = json.Unmarshal(out.Item, &promotion)
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		promotions = append(promotions, promotion)
	}

	return promotions, nil
}

func (r *promotionService) GetCategories(ctx context.Context) ([]model.Category, error) {
	out, err := r.db.ScanItem(ctx,
		&aws.ScanItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.category"),
		})

	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	if out == nil || out.Items == nil || len(out.Items) == 0 {
		return nil, nil
	}

	var categories []model.Category

	for _, out := range out.Items {
		category := model.Category{}
		err = json.Unmarshal(out.Item, &category)
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *promotionService) validPromotion(ctx context.Context, promotion *model.Promotion) error {
	if promotion == nil {
		return errors.New("user is nil")
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

	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     promotion.UserID,
					ValueType: aws.String},
			},
		})

	if out == nil || out.Item == nil {
		err = errors.New("user not found")
		r.log.Error(err.Error())
		return err
	}

	return nil
}
