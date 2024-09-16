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
	"regexp"
	"strings"
	"time"
)

func NewUserService(
	db aws.DynamoDB,
	s3 aws.BucketS3,
	cfg *config.Config,
	log config.Logger,
) handler.UserHandler {
	return &userService{
		db:  db,
		s3:  s3,
		cfg: cfg,
		log: log,
	}
}

type userService struct {
	db  aws.DynamoDB
	s3  aws.BucketS3
	cfg *config.Config
	log config.Logger
}

func (r *userService) CreateUser(ctx context.Context, user *model.User) error {

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

func (r *userService) UpdateUserPicture(ctx context.Context, id string, image io.Reader) error {

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

func (r *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
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

func (r *userService) validUser(user *model.User) error {
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

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
