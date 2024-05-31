package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
// Invoice represents an invoice for an order
// @Description Invoice details
// @Description swagger:model
type Invoice struct {
    // ID of the invoice
    // @swag.Type string
    ID primitive.ObjectID `bson:"_id"`

    // Unique ID of the invoice
    Invoice_id string `json:"invoice_id"`

    // ID of the associated order
    Order_id string `json:"order_id"`

    // Payment method used for the order
    // @swag.Required
    Payment_method *string `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`

    // Payment status of the order
    // @swag.Required
    Payment_status *string `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`

    // Due date for the payment
    Payment_due_date time.Time `json:"Payment_due_date"`

    // Created timestamp of the invoice
    Created_at time.Time `json:"created_at"`

    // Updated timestamp of the invoice
    Updated_at time.Time `json:"updated_at"`
}