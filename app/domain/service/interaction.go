package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port/handler"
	"pixelPromo/domain/port/repository"
	"strings"
	"time"
)

func NewInteractionService(
	interactionRp repository.InteractionRepository,
	userRp repository.UserRepository,
	cfg *config.Config,
	log config.Logger,
) handler.InteractionHandler {
	return &interactionService{
		interactionRp: interactionRp,
		userRp:        userRp,
		cfg:                   cfg,
		log:                   log,
	}
}

type interactionService struct {
	interactionRp repository.InteractionRepository,
	userRp repository.UserRepository,
	cfg *config.Config
	log config.Logger
}

func (r *interactionService) GetInteractionByID(ctx context.Context, id string) (model.PromotionInteraction, error) {
	interaction, err := r.interactionRp.GetInteractionByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return model.PromotionInteraction{}, err
	}

	return interaction, nil
}

func (r *interactionService) CreateInteraction(ctx context.Context, interaction model.PromotionInteraction) error {

	err := r.validInteraction(interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	err = r.interactionRp.CreateInteraction(ctx, interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	ownerItem, err := r.userRp.Get(ctx,
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

func (r *interactionService) validInteraction(interaction model.PromotionInteraction) error {

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

func (r *interactionService) CreateUserScoreByInteraction(interaction *model.PromotionInteraction) (*model.UserScore, error) {

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

func (r *interactionService) getPointsByInteractionType(interactionType model.InteractionType) (int, error) {
	points := r.cfg.Viper.GetInt(fmt.Sprintf("service.score.interactions.%s", interactionType))
	if points <= 0 {
		return 0, errors.New("interaction points not found")
	}
	return points, nil
}

func (r *interactionService) EditUserStatisticByScore(ctx context.Context, user *model.User, score *model.UserScore) (*model.User, error) {

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

func (r *interactionService) calculateElo(ctx context.Context, user *model.User, newPoints int) (string, error) {

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
