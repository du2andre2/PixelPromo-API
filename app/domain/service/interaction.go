package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"pixelPromo/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"strings"
	"time"
)

func NewInteractionService(
	rp port.Repository,
	cfg *config.Config,
	log config.Logger,
) port.InteractionHandler {
	return &interactionService{
		rp:  rp,
		cfg: cfg,
		log: log,
	}
}

type interactionService struct {
	rp  port.Repository
	cfg *config.Config
	log config.Logger
}

func (r *interactionService) GetInteractionByID(ctx context.Context, id string) (*model.PromotionInteraction, error) {
	interaction, err := r.rp.GetInteractionByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return interaction, nil
}

func (r *interactionService) GetCommentsByPromotionID(ctx context.Context, id string) ([]model.PromotionInteraction, error) {
	interaction, err := r.rp.GetCommentsByPromotionID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return interaction, nil
}

func (r *interactionService) GetInteractionsCountersByPromotionID(ctx context.Context, id string) (map[string]int, error) {
	interactions, err := r.rp.GetInteractionsByPromotionID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	counters := map[string]int{
		"favorite": 0,
		"like":     0,
		"comment":  0,
	}
	for _, interaction := range interactions {
		counters[string(interaction.InteractionType)] += 1
	}

	return counters, nil
}

func (r *interactionService) CreateOrUpdateInteraction(ctx context.Context, interaction *model.PromotionInteraction) error {

	err := r.validInteraction(interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	ownerUser, err := r.rp.GetUserByID(ctx, interaction.OwnerUserID)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	if ownerUser == nil {
		return errors.New("owner user not found")
	}

	interaction.CreatedAt = time.Now()
	interaction.ID = fmt.Sprintf("%d", interaction.CreatedAt.UnixNano())

	score, err := r.CreateUserScoreByInteraction(interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	ownerUser, err = r.editUserStatisticByScore(ctx, ownerUser, score)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	err = r.rp.CreateOrUpdateInteraction(ctx, interaction)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	err = r.rp.CreateOrUpdateUser(ctx, ownerUser)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	err = r.rp.CreateOrUpdateUserScore(ctx, score)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("interaction and score created")
	return nil
}

func (r *interactionService) validInteraction(interaction *model.PromotionInteraction) error {
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

	if len(strings.TrimSpace(string(interaction.InteractionType))) == 0 {
		return errors.New("interactionType is empty")
	}

	if !interaction.IsValidType() {
		return errors.New("type is invalid")
	}

	if interaction.InteractionType == model.Comment {
		if len(strings.TrimSpace(interaction.Comment)) == 0 {
			return errors.New("comment is empty")
		}
	}

	return nil
}

func (r *interactionService) CreateUserScoreByInteraction(interaction *model.PromotionInteraction) (*model.UserScore, error) {
	var score model.UserScore

	score.CreatedAt = time.Now()
	score.ID = fmt.Sprintf("%d", score.CreatedAt.UnixNano())
	score.UserID = interaction.OwnerUserID

	points, err := r.getPointsByInteractionType(interaction.InteractionType)
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

func (r *interactionService) editUserStatisticByScore(ctx context.Context, user *model.User, score *model.UserScore) (*model.User, error) {

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

	initDate := time.Now().Add((24 * 7 * time.Hour) * -1)
	scoreList, err := r.rp.GetAllUserScoreByTimeWithUserId(ctx, user.ID, initDate)
	if err != nil {
		r.log.Error(err.Error())
		return "", err
	}

	var pointsInRange int
	pointsInRange += newPoints

	for _, score := range scoreList {
		pointsInRange += score.Points
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
