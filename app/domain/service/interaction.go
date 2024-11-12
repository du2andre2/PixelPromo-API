package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"pixelPromo/domain/model"
	"strings"
	"time"
)

func (s *service) GetCommentsByPromotionID(ctx context.Context, id string) ([]model.PromotionInteraction, error) {
	interaction, err := s.rp.GetInteractionsByTypeWithPromotionID(ctx, model.Comment, id)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return interaction, nil
}

func (s *service) GetInteractionStatisticsByPromotionID(ctx context.Context, id string) (map[string]int, error) {
	interactions, err := s.rp.GetInteractionsByPromotionID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
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

func (s *service) GetInteractionStatisticsByUserID(ctx context.Context, id string) (map[string]int, error) {
	interactions, err := s.rp.GetInteractionsByUserID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	counters := map[string]int{
		"favorite": 0,
		"like":     0,
		"comment":  0,
		"create":   0,
	}
	for _, interaction := range interactions {
		counters[string(interaction.InteractionType)] += 1
	}

	return counters, nil
}

func (s *service) GetInteractionStatisticsByUserIDWithPromotionID(ctx context.Context, userID string, promotionID string) (map[string]bool, error) {
	interactions, err := s.rp.GetInteractionsByUserIDWithPromotionID(ctx, userID, promotionID)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	counters := map[string]bool{
		"favorite": false,
		"like":     false,
	}
	for _, interaction := range interactions {
		counters[string(interaction.InteractionType)] = true
	}

	return counters, nil
}

func (s *service) CreateInteraction(ctx context.Context, newInteraction *model.PromotionInteraction) error {

	promotion, err := s.rp.GetPromotionByID(ctx, newInteraction.PromotionID)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	if promotion == nil {
		return errors.New("promotion not found")
	}

	newInteraction.OwnerUserID = promotion.UserID
	newInteraction.CreatedAt = time.Now()
	newInteraction.ID = fmt.Sprintf("%s#%s#%s#%s", newInteraction.UserID, newInteraction.OwnerUserID, newInteraction.PromotionID, newInteraction.InteractionType.String())

	if newInteraction.InteractionType == model.Create || newInteraction.InteractionType == model.Comment {
		newInteraction.ID = fmt.Sprintf("%s#%s", newInteraction.ID, newInteraction.CreatedAt.String())
	}

	err = s.validInteraction(newInteraction)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	ownerUser, err := s.rp.GetUserByID(ctx, newInteraction.OwnerUserID)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	if ownerUser == nil {
		return errors.New("owner user not found")
	}

	score, err := s.CreateUserScoreByInteraction(newInteraction)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	interaction, err := s.rp.GetInteractionByID(ctx, newInteraction.ID)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	if interaction != nil && interaction.ID == newInteraction.ID {

		score.Points = score.Points * -1

		ownerUser, err = s.editUserStatisticByScore(ctx, ownerUser, score)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}

		err = s.rp.CreateOrUpdateInteraction(ctx, newInteraction)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}
		err = s.rp.CreateOrUpdateUser(ctx, ownerUser)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}
		err = s.rp.CreateOrUpdateUserScore(ctx, score)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}

		err = s.rp.DeleteInteraction(ctx, newInteraction.ID)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}

		s.log.Debug("interaction and score deleted")
		return nil
	}

	ownerUser, err = s.editUserStatisticByScore(ctx, ownerUser, score)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	err = s.rp.CreateOrUpdateInteraction(ctx, newInteraction)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	err = s.rp.CreateOrUpdateUser(ctx, ownerUser)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	err = s.rp.CreateOrUpdateUserScore(ctx, score)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("interaction and score created")
	return nil
}

func (s *service) validInteraction(interaction *model.PromotionInteraction) error {
	if interaction == nil {
		return errors.New("interaction is nil")
	}

	if len(strings.TrimSpace(interaction.UserID)) == 0 {
		return errors.New("userId is empty")
	}

	if len(strings.TrimSpace(interaction.OwnerUserID)) == 0 {
		return errors.New("ownerUserId is empty")
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

func (s *service) CreateUserScoreByInteraction(interaction *model.PromotionInteraction) (*model.UserScore, error) {
	var score model.UserScore

	score.CreatedAt = time.Now()
	score.ID = fmt.Sprintf("%d", score.CreatedAt.UnixNano())
	score.UserID = interaction.OwnerUserID

	points, err := s.getPointsByInteractionType(interaction.InteractionType)
	if err != nil {
		return nil, err
	}
	score.Points = points

	return &score, nil
}

func (s *service) getPointsByInteractionType(interactionType model.InteractionType) (int, error) {
	points := s.cfg.Viper.GetInt(fmt.Sprintf("service.score.interactions.%s", interactionType))
	if points <= 0 {
		return 0, errors.New("interaction points not found")
	}
	return points, nil
}

func (s *service) editUserStatisticByScore(ctx context.Context, user *model.User, score *model.UserScore) (*model.User, error) {
	minimalPointsLevel := s.cfg.Viper.GetInt("service.score.level.minimalPointsLevel")
	growthRate := s.cfg.Viper.GetFloat64("service.score.level.growthRate")

	newLevel := calculateLevel(user, score.Points, minimalPointsLevel, growthRate)
	newElo, err := s.calculateElo(ctx, user, score.Points)
	if err != nil {
		return nil, err
	}

	user.TotalScore += score.Points
	user.Elo = newElo
	user.Level = newLevel
	return user, nil
}

func (s *service) calculateElo(ctx context.Context, user *model.User, newPoints int) (string, error) {
	initDate := time.Now().AddDate(0, 0, -s.cfg.Viper.GetInt("service.score.elo.timeRangeInDays"))
	scoreList, err := s.rp.GetAllUserScoreByTimeWithUserId(ctx, user.ID, initDate)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	pointsInRange := newPoints
	for _, score := range scoreList {
		pointsInRange += score.Points
	}

	switch {
	case pointsInRange >= s.cfg.Viper.GetInt("service.score.elo.levels.diamond.minimal-score"):
		return "diamond", nil
	case pointsInRange >= s.cfg.Viper.GetInt("service.score.elo.levels.platinum.minimal-score"):
		return "platinum", nil
	case pointsInRange >= s.cfg.Viper.GetInt("service.score.elo.levels.gold.minimal-score"):
		return "gold", nil
	case pointsInRange >= s.cfg.Viper.GetInt("service.score.elo.levels.silver.minimal-score"):
		return "silver", nil
	case pointsInRange >= s.cfg.Viper.GetInt("service.score.elo.levels.bronze.minimal-score"):
		return "bronze", nil
	default:
		return "none", nil
	}
}

func pointsRequiredForLevel(level int, minimalPointsLevel int, growthRate float64) int {
	return int(float64(minimalPointsLevel) * math.Pow(growthRate, float64(level-1)))
}

func calculateLevel(user *model.User, points int, minimalPointsLevel int, growthRate float64) int {

	currentLevel := user.Level
	for points >= pointsRequiredForLevel(currentLevel, minimalPointsLevel, growthRate) {
		points -= pointsRequiredForLevel(currentLevel, minimalPointsLevel, growthRate)
		currentLevel++
	}
	return currentLevel
}
