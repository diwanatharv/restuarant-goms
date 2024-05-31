package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Table represents a table in the restaurant
// @Description Table details
// @Description swagger:model
type Table struct {
	// ID of the table
	// @swag.Type string
	ID primitive.ObjectID `bson:"_id"`

	// Number of guests allowed at the table
	// @swag.Required
	Number_of_guests *int `json:"number_of_guests" validate:"required"`

	// Table number
	// @swag.Required
	Table_number *int `json:"table_number" validate:"required"`

	// Created timestamp of the table
	Created_at time.Time `json:"created_at"`

	// Updated timestamp of the table
	Updated_at time.Time `json:"updated_at"`

	// Unique ID of the table
	Table_id string `json:"table_id"`
}
