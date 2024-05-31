package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Menu represents a menu item
// @Description Menu details
// @Description swagger:model
type Menu struct {
    // ID of the menu item
    // @swag.Type string
    ID primitive.ObjectID `bson:"_id"`

    // Name of the menu item
    // @swag.Required
    Name string `json:"name" validate:"required"`

    // Category of the menu item
    // @swag.Required
    Category string `json:"category" validate:"required"`

    // Start date of the menu item availability
    Start_Date *time.Time `json:"start_date"`

    // End date of the menu item availability
    End_Date *time.Time `json:"end_date"`

    // Created timestamp of the menu item
    Created_at time.Time `json:"created_at"`

    // Updated timestamp of the menu item
    Updated_at time.Time `json:"updated_at"`

    // Unique ID of the menu item
    Menu_id string `json:"food_id"`
}
