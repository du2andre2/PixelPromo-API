package model

import "time"

type User struct {
	ID         string    `json:"id" dynamodbav:"id"` //PK
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	PictureUrl string    `json:"pictureUrl"`
	TotalScore int       `json:"totalScore"`
	Level      string    `json:"level"`
	Elo        string    `json:"elo"`
	CreatedAt  time.Time `json:"createdAt"`
}

type UserScore struct {
	ID        string    `json:"id"` //PK
	UserID    string    `json:"userId"`
	Points    int       `json:"points"`
	ScoreDate time.Time `json:"scoreDate"`
}
