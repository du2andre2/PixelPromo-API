package model

import "time"

type User struct {
	Id         string    `json:"id" dynamodbav:"id"` //PK
	Email      string    `json:"email" dynamodbav:"email"`
	Name       string    `json:"name" dynamodbav:"name"`
	Password   string    `json:"password" dynamodbav:"password"`
	PictureUrl string    `json:"pictureUrl" dynamodbav:"pictureUrl"`
	TotalScore int       `json:"totalScore" dynamodbav:"totalScore"`
	Level      int       `json:"level" dynamodbav:"level"`
	Elo        string    `json:"elo" dynamodbav:"elo"`
	CreatedAt  time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

type Login struct {
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
}

type UserScore struct {
	Id        string    `json:"id" dynamodbav:"id"` //PK
	UserId    string    `json:"userId" dynamodbav:"userId"`
	Points    int       `json:"points" dynamodbav:"points"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
}
