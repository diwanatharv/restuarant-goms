package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents an order
// @Description Order details
// @Description swagger:model
type Order struct {
	// ID of the order
	// @swag.Type string
	ID primitive.ObjectID `bson:"_id"`

	// Date of the order
	// @swag.Required
	Order_Date time.Time `json:"order_date" validate:"required"`

	// Created timestamp of the order
	Created_at time.Time `json:"created_at"`

	// Updated timestamp of the order
	Updated_at time.Time `json:"updated_at"`

	// Unique ID of the order
	Order_id string `json:"order_id"`

	// ID of the associated table
	// @swag.Required
	Table_id *string `json:"table_id" validate:"required"`
}
