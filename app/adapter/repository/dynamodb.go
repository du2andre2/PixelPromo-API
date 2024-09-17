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

func (r repository) GetAllUserScoreByRange(ctx context.Context, userID string, limit int) ([]model.UserScore, error) {
	tableName := r.cfg.Viper.GetString("aws.dynamodb.tables.user-score")

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userID},
		},
		Limit: aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, err
	}

	if result == nil || result.Items == nil {
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

	if query.Search != "" {
		searchExpr := "(contains(title, :search) OR contains(description, :search) OR contains(platform, :search) OR contains(link, :search))"
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

	var scanInput dynamodb.ScanInput
	if len(filterExprs) > 0 {
		scanInput = dynamodb.ScanInput{
			TableName:                 aws.String(tableName),
			FilterExpression:          aws.String(filterExpr),
			ExpressionAttributeValues: exprAttrValues,
		}
	} else {
		scanInput = dynamodb.ScanInput{
			TableName: aws.String(tableName),
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