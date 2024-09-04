package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/model"
	"regexp"
	"strings"
	"time"
)

type Repository interface {
	CreateInteraction(context.Context, *model.PromotionInteraction) error

	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(context.Context, string) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
	GetAllPromotions(context.Context, model.PromotionQuery) ([]model.Promotion, error)
	GetPromotionByCategory(context.Context, string) ([]model.Promotion, error)
	GetCategories(context.Context) ([]model.Category, error)
}

func NewRepository(
	db aws.DynamoDB,
	s3 aws.BucketS3,
	cfg *config.Config,
	log config.Logger,
) Repository {
	return &repository{
		db:  db,
		s3:  s3,
		cfg: cfg,
		log: log,
	}
}

type repository struct {
	db  aws.DynamoDB
	s3  aws.BucketS3
	cfg *config.Config
	log config.Logger
}

func (r *repository) CreateInteraction(ctx context.Context, interaction *model.PromotionInteraction) error {

	err := r.validInteraction(interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	userItem, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     interaction.UserID,
					ValueType: aws.String},
			},
		})

	if userItem == nil || userItem.Item == nil {
		err = errors.New("user not found")
		r.log.Error(err.Error())
		return err
	}

	ownerItem, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     interaction.OwnerUserID,
					ValueType: aws.String},
			},
		})

	if ownerItem == nil || ownerItem.Item == nil {
		err = errors.New("ownerUser not found")
		r.log.Error(err.Error())
		return err
	}

	promotionItem, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     interaction.PromotionID,
					ValueType: aws.String},
			},
		})

	if promotionItem == nil || promotionItem.Item == nil {
		err = errors.New("promotion not found")
		r.log.Error(err.Error())
		return err
	}

	interaction.InteractionDate = time.Now()
	interaction.ID = fmt.Sprintf("%d", interaction.InteractionDate.UnixNano())

	score, err := r.CreateUserScoreByInteraction(interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	owner := &model.User{}
	err = json.Unmarshal(ownerItem.Item, owner)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	owner, err = r.EditUserStatisticByScore(ctx, owner, score)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	if err = r.db.BatchPutItem(ctx,
		[]aws.PutItemInput{
			{
				TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction"),
				BodyItem:  interaction,
			}, {
				TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user-score"),
				BodyItem:  score,
			}, {
				TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
				BodyItem:  owner,
			},
		}); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("interaction and score created")
	return nil
}

func (r *repository) CreateUser(ctx context.Context, user *model.User) error {

	err := r.validUser(user)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	user.CreatedAt = time.Now()
	user.ID = fmt.Sprintf("%d", user.CreatedAt.UnixNano())

	if err = r.db.PutItem(ctx,
		&aws.PutItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			BodyItem:  user,
		}); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("user created")
	return nil
}

func (r *repository) UpdateUserPicture(ctx context.Context, id string, image io.Reader) error {

	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			Keys: []aws.Key{
				{
					Name:      "id",
					Value:     id,
					ValueType: aws.String},
			},
		})

	if out == nil || out.Item == nil {
		err = errors.New("user not found")
		r.log.Error(err.Error())
		return err
	}

	outS3, err := r.s3.PutFile(ctx,
		&aws.PutFileInput{
			BucketName: r.cfg.Viper.GetString("aws.s3.buckets.user-pictures"),
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

	user := &model.User{}
	err = json.Unmarshal(out.Item, user)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	user.PictureUrl = outS3.UrlFile

	if err = r.db.PutItem(ctx,
		&aws.PutItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			BodyItem:  user,
		}); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("picture uploaded")
	return nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
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

	user := &model.User{}
	err = json.Unmarshal(out.Item, user)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (r *repository) CreatePromotion(ctx context.Context, promotion *model.Promotion) error {
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

func (r *repository) UpdatePromotionImage(ctx context.Context, id string, image io.Reader) error {
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

func (r *repository) GetPromotionByID(ctx context.Context, id string) (*model.Promotion, error) {
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

func (r *repository) GetAllPromotions(ctx context.Context, params model.PromotionQuery) ([]model.Promotion, error) {
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

func (r *repository) GetPromotionByCategory(ctx context.Context, category string) ([]model.Promotion, error) {
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

func (r *repository) GetCategories(ctx context.Context) ([]model.Category, error) {
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

func (r *repository) validUser(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if len(strings.TrimSpace(user.Email)) == 0 {
		return errors.New("email is empty")
	}
	if len(strings.TrimSpace(user.Name)) == 0 {
		return errors.New("name is empty")
	}
	if len(strings.TrimSpace(user.Password)) == 0 {
		return errors.New("password is empty")
	}
	if !isEmailValid(user.Email) {
		return errors.New("user email is invalid")
	}
	return nil
}

func (r *repository) validPromotion(ctx context.Context, promotion *model.Promotion) error {
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

func (r *repository) validInteraction(interaction *model.PromotionInteraction) error {
	if interaction == nil {
		return errors.New("interaction is nil")
	}

	if len(strings.TrimSpace(interaction.UserID)) == 0 {
		return errors.New("userId is empty")
	}

	if len(strings.TrimSpace(interaction.OwnerUserID)) == 0 {
		return errors.New("ownerUserId is empty")
	}

	if interaction.OwnerUserID == interaction.UserID {
		return errors.New("owner cannot interact with the promotion")
	}

	if len(strings.TrimSpace(interaction.PromotionID)) == 0 {
		return errors.New("promotionId is empty")
	}

	if len(strings.TrimSpace(string(interaction.Type))) == 0 {
		return errors.New("type is empty")
	}

	if !interaction.IsValidType() {
		return errors.New("type is invalid")
	}

	if interaction.Type == model.Comment {
		if len(strings.TrimSpace(interaction.Comment)) == 0 {
			return errors.New("comment is empty")
		}
	}

	return nil
}

func (r *repository) CreateUserScoreByInteraction(interaction *model.PromotionInteraction) (*model.UserScore, error) {

	var score model.UserScore

	score.ScoreDate = time.Now()
	score.ID = fmt.Sprintf("%d", score.ScoreDate.UnixNano())
	score.UserID = interaction.OwnerUserID

	points, err := r.getPointsByInteractionType(interaction.Type)
	if err != nil {
		return nil, err
	}
	score.Points = points

	return &score, nil
}

func (r *repository) getPointsByInteractionType(interactionType model.InteractionType) (int, error) {
	points := r.cfg.Viper.GetInt(fmt.Sprintf("service.score.interactions.%s", interactionType))
	if points <= 0 {
		return 0, errors.New("interaction points not found")
	}
	return points, nil
}

func (r *repository) EditUserStatisticByScore(ctx context.Context, user *model.User, score *model.UserScore) (*model.User, error) {

	newLevel := calculateLevel(user, score.Points)
	newElo, err := r.calculateElo(ctx, user, score.Points)
	if err != nil {
		return nil, err
	}

	user.TotalScore += score.Points
	user.Elo = newElo
	user.Level = newLevel
	return user, nil
}

func (r *repository) calculateElo(ctx context.Context, user *model.User, newPoints int) (string, error) {

	rangeInDays := -7
	dateRange := time.Now().AddDate(0, 0, rangeInDays)
	out, err := r.db.ScanItem(ctx,
		&aws.ScanItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user-score"),
			Conditions: []aws.ConditionParam{
				{
					Names:         []string{"userId"},
					Value:         user.ID,
					OperationType: aws.Equal,
				},
				{
					Names:         []string{"scoreDate"},
					Value:         dateRange,
					OperationType: aws.GreaterThanEqual,
				},
			},
		})

	if err != nil {
		r.log.Error(err.Error())
		return "", err
	}

	var pointsInRange int
	pointsInRange += newPoints

	if !(out == nil || out.Items == nil || len(out.Items) == 0) {
		for _, out := range out.Items {
			userScore := model.UserScore{}
			err = json.Unmarshal(out.Item, &userScore)
			if err != nil {
				r.log.Error(err.Error())
				return "", err
			}
			pointsInRange += userScore.Points
		}
	}

	if pointsInRange >= 100 { //todo:refactor this
		return "silver", nil
	} else if pointsInRange >= 25 {
		return "bronze", nil
	} else {
		return "none", nil
	}
}

func pointsRequiredForLevel(level int) int {
	minimalPointsLevel := 10
	growthRate := 1.30
	return int(float64(minimalPointsLevel) * math.Pow(growthRate, float64(level-1)))
}

func calculateLevel(user *model.User, points int) int {
	currentLevel := user.Level //todo: refactor this
	for points >= pointsRequiredForLevel(currentLevel) {
		points -= pointsRequiredForLevel(currentLevel)
		currentLevel++
	}
	return currentLevel
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
