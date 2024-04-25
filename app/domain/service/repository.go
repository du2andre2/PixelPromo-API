package service

import (
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/model"
)

type Repository interface {
	GetUser() (*model.User, error)
	GetUserByID(id string) (*model.User, error)
}

type repository struct {
	db  aws.DynamoDB
	s3  aws.BucketS3
	cfg *config.Config
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

func (r *repository) GetUser() (*model.User, error) {
	_ = r.db.GetItem()
	_ = r.s3.GetFile()
	return &model.User{
		ID: "1",
	}, nil
}

func (r *repository) GetUserByID(id string) (*model.User, error) {
	if id == "1" {
		return &model.User{ID: "1"}, nil
	}
	return nil, nil
}
