package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note represents a note
// @Description Note details
// @Description swagger:model
type Note struct {
	// ID of the note
	// @swag.Type string
	ID primitive.ObjectID `bson:"_id"`

	// Text content of the note
	Text string `json:"text"`

	// Title of the note
	Title string `json:"title"`

	// Created timestamp of the note
	Created_at time.Time `json:"created_at"`

	// Updated timestamp of the note
	Updated_at time.Time `json:"updated_at"`

	// Unique ID of the note
	Note_id string `json:"note_id"`
}
