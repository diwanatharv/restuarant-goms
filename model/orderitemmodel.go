package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents an item in an order
// @Description Order item details
// @Description swagger:model
type OrderItem struct {
	// ID of the order item
	// @swag.Type string
	ID primitive.ObjectID `bson:"_id"`

	// Quantity of the order item
	// @swag.Required
	Quantity *string `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`

	// Unit price of the order item
	// @swag.Required
	Unit_price *float64 `json:"unit_price" validate:"required"`

	// Created timestamp of the order item
	Created_at time.Time `json:"created_at"`

	// Updated timestamp of the order item
	Updated_at time.Time `json:"updated_at"`

	// ID of the associated food item
	// @swag.Required
	Food_id *string `json:"food_id" validate:"required"`

	// Unique ID of the order item
	Order_item_id string `json:"order_item_id"`

	// ID of the associated order
	// @swag.Required
	Order_id string `json:"order_id" validate:"required"`
}
