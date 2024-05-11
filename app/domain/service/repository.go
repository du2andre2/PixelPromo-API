package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/model"
)

type Repository interface {
	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(string) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, io.Reader) error
	GetPromotionByID(string) (*model.Promotion, error)
}

func NewRepository(
	db aws.DynamoDB,
	s3 aws.BucketS3,
	cfg *config.Config,
) Repository {
	return &repository{
		db:  db,
		s3:  s3,
		cfg: cfg,
	}
}

type repository struct {
	db  aws.DynamoDB
	s3  aws.BucketS3
	cfg *config.Config
}

func (r repository) CreateUser(ctx context.Context, user *model.User) error {

	err := r.db.PutItem(ctx,
		&aws.PutItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			BodyItem:  user,
		})

	return err
}

func (r repository) UpdateUserPicture(ctx context.Context, id string, reader io.Reader) error {

	out, err := r.db.GetItem(ctx,
		&aws.GetItemInput{
			TableName: r.cfg.Viper.GetString("aws.dynamodb.tables.user"),
			Keys: []aws.Key{
				{Name: "id", Value: id, ValueType: aws.String},
			},
		})

	if out == nil || out.Item == nil {
		return errors.New("user not found")
	}

	r.s3.PutFile(ctx,
		&aws.PutFileInput{
			BucketName: r.cfg.Viper.GetString("aws.s3.buckets.user-pictures"),
			FileName:   fmt.Sprintf("%s.jpg", id),
			BodyItem:   reader,
		},
	)

	return err
}

func (r repository) GetUserByID(s string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) CreatePromotion(ctx context.Context, promotion *model.Promotion) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) UpdatePromotionImage(ctx context.Context, reader io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetPromotionByID(s string) (*model.Promotion, error) {
	//TODO implement me
	panic("implement me")
}
