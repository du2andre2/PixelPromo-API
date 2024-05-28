package model

import "time"

type User struct {
	ID         string    `json:"id" dynamodbav:"id"` //PK
	Email      string    `json:"email" dynamodbav:"email"`
	Name       string    `json:"name" dynamodbav:"name"`
	Password   string    `json:"password" dynamodbav:"password"`
	PictureUrl string    `json:"pictureUrl" dynamodbav:"pictureUrl"`
	TotalScore int       `json:"totalScore" dynamodbav:"totalScore"`
	Level      int       `json:"level" dynamodbav:"level"`
	Elo        string    `json:"elo" dynamodbav:"elo"`
	CreatedAt  time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

type UserScore struct {
	ID        string    `json:"id" dynamodbav:"id"` //PK
	UserID    string    `json:"userId" dynamodbav:"userId"`
	Points    int       `json:"points" dynamodbav:"points"`
	ScoreDate time.Time `json:"scoreDate" dynamodbav:"scoreDate"`
}
