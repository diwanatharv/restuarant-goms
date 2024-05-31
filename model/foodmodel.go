package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
// Food represents a food item
// @Description Food item details
// @Description swagger:model
type Food struct {
    // ID of the food item
    // @swag.Type string
    ID primitive.ObjectID `bson:"_id"`

    // Name of the food item
    // @swag.Required
    Name *string `json:"name" validate:"required,min=2,max=100"`

    // Price of the food item
    // @swag.Required
    Price *float64 `json:"price" validate:"required"`

    // Image URL of the food item
    // @swag.Required
    Food_image *string `json:"food_image" validate:"required"`

    // Created timestamp of the food item
    Created_at time.Time `json:"created_at"`

    // Updated timestamp of the food item
    Updated_at time.Time `json:"updated_at"`

    // Unique ID of the food item
    Food_id string `json:"food_id"`

    // ID of the menu associated with the food item
    // @swag.Required
    Menu_id *string `json:"menu_id" validate:"required"`
}