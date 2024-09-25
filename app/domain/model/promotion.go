package model

import "time"

type Promotion struct {
	ID              string    `json:"id" dynamodbav:"id"` //PK
	UserID          string    `json:"userId" dynamodbav:"userId"`
	Title           string    `json:"title" dynamodbav:"title"`
	Description     string    `json:"description" dynamodbav:"description"`
	OriginalPrice   string    `json:"originalPrice" dynamodbav:"originalPrice"`
	DiscountedPrice string    `json:"discountedPrice" dynamodbav:"discountedPrice"`
	DiscountBadge   string    `json:"discountBadge" dynamodbav:"discountBadge"`
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
	ID              string          `json:"id" dynamodbav:"id"` //PK
	PromotionID     string          `json:"promotionId" dynamodbav:"promotionId"`
	OwnerUserID     string          `json:"ownerUserId" dynamodbav:"ownerUserId"`
	UserID          string          `json:"userId" dynamodbav:"userId"`
	Comment         string          `json:"comment" dynamodbav:"comment"`
	InteractionType InteractionType `json:"interactionType" dynamodbav:"interactionType"`
	InteractionDate time.Time       `json:"interactionDate" dynamodbav:"interactionDate"`
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

const (
	Favorite InteractionType = "favorite"
	Like     InteractionType = "like"
	Comment  InteractionType = "comment"
	Create   InteractionType = "create"
)

type PromotionQuery struct {
	Categories []string `json:"category"`
	Search     string   `json:"search"`
}
