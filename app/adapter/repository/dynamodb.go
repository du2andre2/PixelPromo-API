package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"pixelPromo/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"strings"
	"time"
)

func NewDynamoDBRepository(
	awsCfg *aws.Config,
	cfg *config.Config,
) port.Repository {
	return &repository{
		client: dynamodb.NewFromConfig(*awsCfg),
		cfg:    cfg,
	}
}

type repository struct {
	client *dynamodb.Client
	cfg    *config.Config
}

func (r repository) CreateOrUpdateInteraction(ctx context.Context, interaction *model.PromotionInteraction) error {
	item, err := attributevalue.MarshalMap(interaction)
	if err != nil {
		return err
	}

	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func (r repository) DeleteInteraction(ctx context.Context, id string) error {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err

}

func (r repository) GetInteractionByID(ctx context.Context, id string) (*model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Item == nil {
		return nil, nil
	}

	var interaction model.PromotionInteraction
	err = attributevalue.UnmarshalMap(result.Item, &interaction)
	if err != nil {
		return nil, err
	}

	return &interaction, nil
}

func (r repository) GetInteractionsByPromotionID(ctx context.Context, id string) ([]model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("promotionId = :promotionId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":promotionId": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var interactions []model.PromotionInteraction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &interactions)
	if err != nil {
		return nil, err
	}

	return interactions, nil
}

func (r repository) GetInteractionsByUserID(ctx context.Context, id string) ([]model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("ownerUserId = :ownerUserId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ownerUserId": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var interactions []model.PromotionInteraction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &interactions)
	if err != nil {
		return nil, err
	}

	return interactions, nil
}
func (r repository) GetInteractionsByUserIDWithPromotionID(ctx context.Context, userID string, promotionID string) ([]model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("promotionId = :promotionId AND userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":      &types.AttributeValueMemberS{Value: userID},
			":promotionId": &types.AttributeValueMemberS{Value: promotionID},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var interactions []model.PromotionInteraction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &interactions)
	if err != nil {
		return nil, err
	}

	return interactions, nil
}

func (r repository) GetInteractionsByTypeWithPromotionID(ctx context.Context, interactionType model.InteractionType, id string) ([]model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("promotionId = :promotionId AND interactionType = :interactionType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":promotionId":     &types.AttributeValueMemberS{Value: id},
			":interactionType": &types.AttributeValueMemberS{Value: interactionType.String()},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var interactions []model.PromotionInteraction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &interactions)
	if err != nil {
		return nil, err
	}

	return interactions, nil
}

func (r repository) GetInteractionsByTypeWithUserID(ctx context.Context, interactionType model.InteractionType, id string) ([]model.PromotionInteraction, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion-interaction")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("userId = :userId AND interactionType = :interactionType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":          &types.AttributeValueMemberS{Value: id},
			":interactionType": &types.AttributeValueMemberS{Value: interactionType.String()},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var interactions []model.PromotionInteraction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &interactions)
	if err != nil {
		return nil, err
	}

	return interactions, nil
}

func (r repository) CreateOrUpdateUser(ctx context.Context, user *model.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user")
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	return err
}

func (r repository) CreateOrUpdateUserScore(ctx context.Context, score *model.UserScore) error {
	item, err := attributevalue.MarshalMap(score)
	if err != nil {
		return err
	}

	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user-score")
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func (r repository) GetAllUserScoreByTimeWithUserId(ctx context.Context, userID string, createdAt time.Time) ([]model.UserScore, error) {
	createdAtISO := createdAt.Format(time.RFC3339)
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user-score")

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		IndexName:        aws.String("CreatedAtIndex"),
		FilterExpression: aws.String("userId = :userId AND createdAt > :createdAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":    &types.AttributeValueMemberS{Value: userID},
			":createdAt": &types.AttributeValueMemberS{Value: createdAtISO},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan table %s using index %s: %w", tableName, "CreatedAtIndex", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var scores []model.UserScore
	err = attributevalue.UnmarshalListOfMaps(result.Items, &scores)
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (r repository) GetAllUserScoreByTime(ctx context.Context, createdAt time.Time) ([]model.UserScore, error) {
	createdAtISO := createdAt.Format(time.RFC3339)
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user-score")

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		IndexName:        aws.String("CreatedAtIndex"),
		FilterExpression: aws.String("createdAt > :createdAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":createdAt": &types.AttributeValueMemberS{Value: createdAtISO},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan table %s using index %s: %w", tableName, "CreatedAtIndex", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var scores []model.UserScore
	err = attributevalue.UnmarshalListOfMaps(result.Items, &scores)
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (r repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user")
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Item == nil {
		return nil, nil
	}
	var user model.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r repository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil || len(result.Items) == 0 {
		return nil, nil
	}

	var users []model.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (r repository) GetUserByEmailAndPassword(ctx context.Context, email string, password string) (*model.User, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user")
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("email = :email AND password = :password"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email":    &types.AttributeValueMemberS{Value: email},
			":password": &types.AttributeValueMemberS{Value: password},
		},
	})

	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil || len(result.Items) == 0 {
		return nil, nil
	}

	var users []model.User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}
func (r repository) CreateOrUpdatePromotion(ctx context.Context, promotion *model.Promotion) error {
	item, err := attributevalue.MarshalMap(promotion)
	if err != nil {
		return err
	}

	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion")
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func (r repository) GetPromotionByID(ctx context.Context, id string) (*model.Promotion, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion")
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Item == nil {
		return nil, nil
	}
	var promotion model.Promotion
	err = attributevalue.UnmarshalMap(result.Item, &promotion)
	if err != nil {
		return nil, err
	}

	return &promotion, nil
}

func (r repository) GetPromotionsWithParams(ctx context.Context, query *model.PromotionQuery) ([]model.Promotion, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion")

	var filterExprs []string
	exprAttrValues := map[string]types.AttributeValue{}

	if query.UserID != "" {
		userIDExpr := "userId = :userId"
		filterExprs = append(filterExprs, userIDExpr)
		exprAttrValues[":userId"] = &types.AttributeValueMemberS{Value: query.UserID}
	}

	if query.Search != "" {
		searchExpr := " (contains(title, :search) OR contains(description, :search) OR contains(platform, :search) OR contains(link, :search))"
		filterExprs = append(filterExprs, searchExpr)
		exprAttrValues[":search"] = &types.AttributeValueMemberS{Value: query.Search}
	}

	for i, category := range query.Categories {
		if category != "" {
			key := fmt.Sprintf(":category%d", i)
			categoryExpr := fmt.Sprintf("contains(categories, %s)", key)
			filterExprs = append(filterExprs, categoryExpr)
			exprAttrValues[key] = &types.AttributeValueMemberS{Value: category}
		}
	}

	filterExpr := ""
	if len(filterExprs) > 0 {
		filterExpr = strings.Join(filterExprs, " AND ")
	}

	limit := &query.Limit
	if query.Limit == 0 {
		limit = nil
	}

	var scanInput dynamodb.ScanInput
	if len(filterExprs) > 0 {
		scanInput = dynamodb.ScanInput{
			TableName:                 aws.String(tableName),
			FilterExpression:          aws.String(filterExpr),
			ExpressionAttributeValues: exprAttrValues,
			Limit:                     limit,
		}
	} else {
		scanInput = dynamodb.ScanInput{
			TableName: aws.String(tableName),
			Limit:     limit,
		}
	}

	result, err := r.client.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var promotions []model.Promotion
	err = attributevalue.UnmarshalListOfMaps(result.Items, &promotions)
	if err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r repository) GetPromotionsByCategory(ctx context.Context, category string) ([]model.Promotion, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.promotion")

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("category = :category"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":category": &types.AttributeValueMemberS{Value: category},
		},
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var promotions []model.Promotion
	err = attributevalue.UnmarshalListOfMaps(result.Items, &promotions)
	if err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r repository) GetCategories(ctx context.Context) ([]model.Category, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.category")

	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, nil
	}
	var categories []model.Category
	err = attributevalue.UnmarshalListOfMaps(result.Items, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
