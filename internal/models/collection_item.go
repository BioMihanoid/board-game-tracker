package models

import "time"

type CollectionItem struct {
	ID            string
	UserID        string
	GameID        string
	Status        CollectionStatus
	Condition     ItemCondition
	PurchasePrice *int
	SoldPrice     *int
	Notes         string
	PurchasedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
