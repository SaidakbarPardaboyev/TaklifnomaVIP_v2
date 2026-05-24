package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseEntity struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type AccountInfo struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

type Organization struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type CurrencyAmount struct {
	Amount   float64  `json:"amount" bson:"amount"`
	Currency Currency `json:"currency" bson:"currency"`
}

type Currency struct {
	ID         int64  `json:"id" bson:"id"`
	IsNational bool   `json:"is_national" bson:"is_national"`
	Name       string `json:"name" bson:"name"`
}

type CurrencyAmountMovement struct {
	CurrencyAmount `bson:",inline"`
	PaymentType    byte `json:"payment_type" bson:"payment_type"`
}
