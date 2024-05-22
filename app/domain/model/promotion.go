package model

import "time"

type Promotion struct {
	ID          string    `json:"id" dynamodbav:"id"` //PK
	UserID      string    `json:"userId" dynamodbav:"userId"`
	Title       string    `json:"title" dynamodbav:"title"`
	Description string    `json:"description" dynamodbav:"description"`
	ImageUrl    string    `json:"imageUrl" dynamodbav:"imageUrl"`
	Link        string    `json:"link" dynamodbav:"link"`
	CategoryIDs []string  `json:"categoryIds" dynamodbav:"categoryIds"`
	CreatedAt   time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

type Category struct {
	ID   string `json:"id" dynamodbav:"id"` //PK
	Name string `json:"name" dynamodbav:"name"`
}

type PromotionInteraction struct {
	ID              string          `json:"id" dynamodbav:"id"` //PK
	PromotionID     string          `json:"promotionId" dynamodbav:"promotionId"`
	UserID          string          `json:"userId" dynamodbav:"userId"`
	Comment         string          `json:"comment" dynamodbav:"comment"`
	Type            InteractionType `json:"type" dynamodbav:"type"`
	InteractionDate time.Time       `json:"interactionDate" dynamodbav:"interactionDate"`
}

type InteractionType string

const (
	Favorite InteractionType = "favorited"
	Like     InteractionType = "like"
	Comment  InteractionType = "comment"
	Create   InteractionType = "create"
)
