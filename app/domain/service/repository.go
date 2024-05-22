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
	"regexp"
	"strings"
	"time"
)

type Repository interface {
	CreateUser(context.Context, *model.User) error
	UpdateUserPicture(context.Context, string, io.Reader) error
	GetUserByID(context.Context, string) (*model.User, error)

	CreatePromotion(context.Context, *model.Promotion) error
	UpdatePromotionImage(context.Context, string, io.Reader) error
	GetPromotionByID(context.Context, string) (*model.Promotion, error)
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

	user := &model.Promotion{}
	err = json.Unmarshal(out.Item, user)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return user, nil
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

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
