package aws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB interface {
	GetItem(context.Context, *GetItemInput) (*GetItemOutput, error)
	PutItem(context.Context, *PutItemInput) error
	DeleteItem(context.Context, *DeleteItemInput) error
	UpdateItem(context.Context, *UpdateItemInput) error
	BatchGetItem(context.Context, *BatchGetItemInput) (*BatchGetItemOutput, error)
	BatchPutItem(context.Context, *BatchPutItemInput) error
	QueryItem(context.Context, *QueryItemInput) (*QueryItemOutput, error)
	ScanItem(context.Context, *ScanItemInput) (*ScanItemOutput, error)
}

type dynamoDB struct {
	client *dynamodb.Client
}

func NewDynamoDb(cfg *aws.Config) DynamoDB {
	return &dynamoDB{
		client: dynamodb.NewFromConfig(*cfg),
	}
}

type ValueType int

const (
	String ValueType = iota
	Int
	Bool
)

type Key struct {
	Name      string
	Value     interface{}
	ValueType ValueType
}

type GetItemInput struct {
	TableName string
	Keys      []Key
}

type GetItemOutput struct {
	Item []byte
}

type PutItemInput struct {
	TableName string
	BodyItem  interface{}
}

type DeleteItemInput struct {
	TableName string
	Keys      []Key
}

type UpdateItemInput struct {
	TableName  string
	Keys       []Key
	Conditions []ConditionParam
	Actions    []ActionParam
}

type ActionParam struct {
	Name       string
	Value      interface{}
	ActionType ActionType
}

type ActionType int

const (
	Set ActionType = iota
)

type ConditionParam struct {
	Name          string
	Value         interface{}
	OperationType OperationType
}

type OperationType int

const (
	Equal OperationType = iota
	Contains
	BeginsWith
)

type BatchGetItemInput struct {
	Inputs []GetItemInput
}

type BatchGetItemOutput struct {
	Items []GetItemOutput
}

type BatchPutItemInput struct {
	TableName string
	Items     []interface{}
}

type QueryItemInput struct {
	TableName  string
	Conditions []ConditionParam
}

type QueryItemOutput struct {
	TableName string
	Items     []GetItemOutput
}

type ScanItemInput struct {
	TableName  string
	Conditions []ConditionParam
}

type ScanItemOutput struct {
	TableName string
	Items     []GetItemOutput
}

func (d dynamoDB) GetItem(ctx context.Context, input *GetItemInput) (*GetItemOutput, error) {
	keys := make(map[string]types.AttributeValue)

	for _, i := range input.Keys {
		switch i.ValueType {
		case String:
			keys[i.Name] = &types.AttributeValueMemberS{Value: i.Value.(string)}
		case Int:
			keys[i.Name] = &types.AttributeValueMemberN{Value: i.Value.(string)}
		case Bool:
			keys[i.Name] = &types.AttributeValueMemberBOOL{Value: i.Value.(bool)}
		}
	}

	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            keys,
		TableName:      aws.String(input.TableName),
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	var output interface{}
	if len(result.Item) > 0 {
		if err := attributevalue.UnmarshalMap(result.Item, &output); err != nil {
			return nil, err
		}
	}
	byteOutput, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return &GetItemOutput{
		Item: byteOutput,
	}, nil
}

func (d dynamoDB) PutItem(ctx context.Context, input *PutItemInput) error {
	item, err := attributevalue.MarshalMap(input.BodyItem)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(input.TableName),
		Item:      item,
	})
	return err
}

func (d dynamoDB) DeleteItem(ctx context.Context, input *DeleteItemInput) error {
	keys := make(map[string]types.AttributeValue)

	for _, i := range input.Keys {
		switch i.ValueType {
		case String:
			keys[i.Name] = &types.AttributeValueMemberS{Value: i.Value.(string)}
		case Int:
			keys[i.Name] = &types.AttributeValueMemberN{Value: i.Value.(string)}
		case Bool:
			keys[i.Name] = &types.AttributeValueMemberBOOL{Value: i.Value.(bool)}
		}
	}

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       keys,
		TableName: aws.String(input.TableName),
	})

	return err
}

func (d dynamoDB) UpdateItem(ctx context.Context, input *UpdateItemInput) error {
	keys := make(map[string]types.AttributeValue)

	for _, i := range input.Keys {
		switch i.ValueType {
		case String:
			keys[i.Name] = &types.AttributeValueMemberS{Value: i.Value.(string)}
		case Int:
			keys[i.Name] = &types.AttributeValueMemberN{Value: i.Value.(string)}
		case Bool:
			keys[i.Name] = &types.AttributeValueMemberBOOL{Value: i.Value.(bool)}
		}
	}

	expr, err := d.buildUpdateExpression(input.Actions, input.Conditions)

	_, err = d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       keys,
		TableName:                 aws.String(input.TableName),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		UpdateExpression:          expr.Update(),
	})

	return err
}

func (d dynamoDB) BatchGetItem(ctx context.Context, input *BatchGetItemInput) (*BatchGetItemOutput, error) {

	requestItems := make(map[string]types.KeysAndAttributes)
	for _, itemInput := range input.Inputs {
		keys := make(map[string]types.AttributeValue)
		for _, i := range itemInput.Keys {
			switch i.ValueType {
			case String:
				keys[i.Name] = &types.AttributeValueMemberS{Value: i.Value.(string)}
			case Int:
				keys[i.Name] = &types.AttributeValueMemberN{Value: i.Value.(string)}
			case Bool:
				keys[i.Name] = &types.AttributeValueMemberBOOL{Value: i.Value.(bool)}
			}
		}

		if _, exist := requestItems[itemInput.TableName]; exist {
			requestItems[itemInput.TableName] = types.KeysAndAttributes{
				Keys:           append(requestItems[itemInput.TableName].Keys),
				ConsistentRead: aws.Bool(true),
			}
		}

		requestItems[itemInput.TableName] = types.KeysAndAttributes{
			Keys: []map[string]types.AttributeValue{
				keys,
			},
			ConsistentRead: aws.Bool(true),
		}

	}

	result, err := d.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	})

	if err != nil {
		return nil, err
	}

	response, ok := result.Responses[""]
	if !ok {
		return nil, nil
	}

	var outputsInterface []interface{}
	if err := attributevalue.UnmarshalListOfMaps(response, outputsInterface); err != nil {
		return nil, err
	}

	output := &BatchGetItemOutput{
		Items: make([]GetItemOutput, 0),
	}

	for _, outputInterface := range outputsInterface {
		byteOutput, err := json.Marshal(outputInterface)
		if err != nil {
			return nil, err
		}

		output.Items = append(output.Items, GetItemOutput{Item: byteOutput})
	}

	return output, nil
}

func (d dynamoDB) BatchPutItem(ctx context.Context, input *BatchPutItemInput) error {
	inputAws := dynamodb.TransactWriteItemsInput{}

	for _, param := range input.Items {
		item, err := attributevalue.MarshalMap(param)
		if err != nil {
			return err
		}

		inputAws.TransactItems = append(inputAws.TransactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(input.TableName),
				Item:      item,
			},
		})
	}

	_, err := d.client.TransactWriteItems(ctx, &inputAws)
	return err
}

func (d dynamoDB) QueryItem(ctx context.Context, input *QueryItemInput) (*QueryItemOutput, error) {
	expr, err := d.buildQueryExpression(input.Conditions)
	outputQuery, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(input.TableName),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.Filter(),
	})

	if err != nil {
		return nil, err
	}

	var outputsInterface []interface{}
	if err := attributevalue.UnmarshalListOfMaps(outputQuery.Items, outputsInterface); err != nil {
		return nil, err
	}

	output := &QueryItemOutput{
		Items: make([]GetItemOutput, 0),
	}

	for _, outputInterface := range outputsInterface {
		byteOutput, err := json.Marshal(outputInterface)
		if err != nil {
			return nil, err
		}

		output.Items = append(output.Items, GetItemOutput{Item: byteOutput})
	}

	return output, nil

}

func (d dynamoDB) ScanItem(ctx context.Context, input *ScanItemInput) (*ScanItemOutput, error) {
	expr, err := d.buildQueryExpression(input.Conditions)
	outputScan, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(input.TableName),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		FilterExpression:          expr.Filter(),
	})

	if err != nil {
		return nil, err
	}

	var outputsInterface []interface{}
	if err := attributevalue.UnmarshalListOfMaps(outputScan.Items, outputsInterface); err != nil {
		return nil, err
	}

	output := &ScanItemOutput{
		Items: make([]GetItemOutput, 0),
	}

	for _, outputInterface := range outputsInterface {
		byteOutput, err := json.Marshal(outputInterface)
		if err != nil {
			return nil, err
		}

		output.Items = append(output.Items, GetItemOutput{Item: byteOutput})
	}

	return output, nil

}

func (d dynamoDB) buildQueryExpression(conditions []ConditionParam) (expression.Expression, error) {

	filterExpr := expression.ConditionBuilder{}

	for i, condition := range conditions {
		switch condition.OperationType {
		case Equal:
			if i == 0 {
				filterExpr = expression.Name(condition.Name).Equal(expression.Value(condition.Value))
			}
			filterExpr = filterExpr.And(expression.Name(condition.Name).Equal(expression.Value(condition.Value)))
		case BeginsWith:
			if i == 0 {
				filterExpr = expression.Name(condition.Name).BeginsWith(condition.Value.(string))
			}
			filterExpr = filterExpr.And(expression.Name(condition.Name).BeginsWith(condition.Value.(string)))
		case Contains:
			if i == 0 {
				filterExpr = expression.Name(condition.Name).Contains(condition.Value.(string))
			}
			filterExpr = filterExpr.And(expression.Name(condition.Name).Contains(condition.Value.(string)))
		}
	}

	return expression.NewBuilder().WithFilter(filterExpr).Build()

}

func (d dynamoDB) buildUpdateExpression(actions []ActionParam, conditions []ConditionParam) (expression.Expression, error) {
	updateExpr := expression.UpdateBuilder{}

	for i, action := range actions {
		switch action.ActionType {
		case Set:
			if i == 0 {
				updateExpr = expression.Set(expression.Name(action.Name), expression.Value(action.Value))
			}
			updateExpr = updateExpr.Set(expression.Name(action.Name), expression.Value(action.Value))
		}
	}
	conditionExpr := expression.ConditionBuilder{}

	for i, condition := range conditions {
		switch condition.OperationType {
		case Equal:
			if i == 0 {
				conditionExpr = expression.Name(condition.Name).Equal(expression.Value(condition.Value))
			}
			conditionExpr = conditionExpr.And(expression.Name(condition.Name).Equal(expression.Value(condition.Value)))
		case BeginsWith:
			if i == 0 {
				conditionExpr = expression.Name(condition.Name).BeginsWith(condition.Value.(string))
			}
			conditionExpr = conditionExpr.And(expression.Name(condition.Name).BeginsWith(condition.Value.(string)))
		case Contains:
			if i == 0 {
				conditionExpr = expression.Name(condition.Name).Contains(condition.Value.(string))
			}
			conditionExpr = conditionExpr.And(expression.Name(condition.Name).Contains(condition.Value.(string)))
		}
	}

	if len(conditions) > 0 {
		return expression.NewBuilder().WithUpdate(updateExpr).WithCondition(conditionExpr).Build()
	}
	return expression.NewBuilder().WithUpdate(updateExpr).Build()

}
