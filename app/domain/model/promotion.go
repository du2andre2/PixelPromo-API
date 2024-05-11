package model

import "time"

type Promotion struct {
	ID          string    `json:"id"` //PK
	UserID      string    `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"imageUrl"`
	Link        string    `json:"link"`
	CategoryIDs []string  `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Category struct {
	ID   string `json:"id"` //PK
	Name string `json:"name"`
}

type PromotionInteraction struct {
	ID              string          `json:"id"` //PK
	PromotionID     string          `json:"promotionId"`
	UserID          string          `json:"userId"`
	Comment         string          `json:"comment"`
	Type            InteractionType `json:"type"`
	InteractionDate time.Time       `json:"interactionDate"`
}

type InteractionType string

const (
	Favorite InteractionType = "favorited"
	Like     InteractionType = "like"
	Comment  InteractionType = "comment"
	Create   InteractionType = "create"
)
