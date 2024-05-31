package model

import (
	 // Load documentation from docs/
	"time"
	
	"github.com/swaggo/swag"
	_ "management/docs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user
// @Description User details
// @Description swagger:model
type User struct {
	// ID of the user
	// @swag.Type string
	ID primitive.ObjectID `bson:"_id"`

	// First name of the user
	// @swag.Required
	First_name *string `json:"first_name" validate:"required,min=2,max=100"`

	// Last name of the user
	// @swag.Required
	Last_name *string `json:"last_name" validate:"required,min=2,max=100"`

	// Password of the user
	// @swag.Required
	Password *string `json:"Password" validate:"required,min=6"`

	// Email of the user
	// @swag.Required
	Email *string `json:"email" validate:"email,required"`

	// Avatar URL of the user
	Avatar *string `json:"avatar"`

	// Phone number of the user
	// @swag.Required
	Phone *string `json:"phone" validate:"required"`

	// Authentication token of the user
	Token *string `json:"token"`

	// Refresh token of the user
	Refresh_Token *string `json:"refresh_token"`

	// Created timestamp of the user
	Created_at time.Time `json:"created_at"`

	// Updated timestamp of the user
	Updated_at time.Time `json:"updated_at"`

	// Unique ID of the user
	User_id string `json:"user_id"`
}
