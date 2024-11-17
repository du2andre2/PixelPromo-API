package model

import "time"

type Promotion struct {
	Id              string    `json:"id" dynamodbav:"id"` //PK
	UserId          string    `json:"userId" dynamodbav:"userId"`
	Title           string    `json:"title" dynamodbav:"title"`
	OriginalPrice   float64   `json:"originalPrice" dynamodbav:"originalPrice"`
	DiscountedPrice float64   `json:"discountedPrice" dynamodbav:"discountedPrice"`
	DiscountBadge   float64   `json:"discountBadge" dynamodbav:"discountBadge"`
	Platform        string    `json:"platform" dynamodbav:"platform"`
	ImageUrl        string    `json:"imageUrl" dynamodbav:"imageUrl"`
	Link            string    `json:"link" dynamodbav:"link"`
	Categories      []string  `json:"categories" dynamodbav:"categories"`
	CreatedAt       time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

type Category struct {
	Name string `json:"name" dynamodbav:"name"`
}

type PromotionInteraction struct {
	Id              string          `json:"id" dynamodbav:"id"` //PK
	PromotionId     string          `json:"promotionId" dynamodbav:"promotionId"`
	OwnerUserId     string          `json:"ownerUserId" dynamodbav:"ownerUserId"`
	UserId          string          `json:"userId" dynamodbav:"userId"`
	Comment         string          `json:"comment" dynamodbav:"comment"`
	InteractionType InteractionType `json:"interactionType" dynamodbav:"interactionType"`
	CreatedAt       time.Time       `json:"createdAt" dynamodbav:"createdAt"`
}

func (p *PromotionInteraction) IsValidType() bool {
	switch p.InteractionType {
	case Create, Comment, Favorite, Like:
		return true
	default:
		return false
	}
}

type InteractionType string

func (t InteractionType) String() string {
	return string(t)
}

const (
	Favorite InteractionType = "favorite"
	Like     InteractionType = "like"
	Comment  InteractionType = "comment"
	Create   InteractionType = "create"
)

type PromotionQuery struct {
	Categories []string `json:"category"`
	Search     string   `json:"search"`
	UserId     string   `json:"userId"`
	Limit      int32    `json:"limit"`
}
