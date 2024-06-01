package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"management/database"
	"management/model"
	"net/http"
	"time"
)

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
// GetInvoices
// @Summary List invoices
// @Description Get all invoices
// @Tags invoices
// @Produce json
// @Success 200 {array} model.CustomBsonM
// @Failure 500 {object} model.ErrorResponse
// @Router /invoices [get]
func GetInvoices(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing invoice items",
		})
	}

	var allInvoices []bson.M
	if err = result.All(ctx, &allInvoices); err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, allInvoices)
}
// GetInvoice
// @Summary Get invoice
// @Description Get an invoice by ID
// @Tags invoices
// @Produce json
// @Param invoice_id path string true "Invoice ID"
// @Success 200 {object} controller.InvoiceViewFormat
// @Failure 500 {object} model.ErrorResponse
// @Router /invoices/{invoice_id} [get]
func GetInvoice(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	invoiceId := c.Param("invoice_id")

	var invoice model.Invoice

	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing invoice item",
		})
	}

	var invoiceView InvoiceViewFormat

	allOrderItems, err := ItemsByOrder(invoice.Order_id)
	invoiceView.Order_id = invoice.Order_id
	invoiceView.Payment_due_date = invoice.Payment_due_date

	invoiceView.Payment_method = "null"
	if invoice.Payment_method != nil {
		invoiceView.Payment_method = *invoice.Payment_method
	}

	invoiceView.Invoice_id = invoice.Invoice_id
	invoiceView.Payment_status = *&invoice.Payment_status
	invoiceView.Payment_due = allOrderItems[0]["payment_due"]
	invoiceView.Table_number = allOrderItems[0]["table_number"]
	invoiceView.Order_details = allOrderItems[0]["order_items"]

	return c.JSON(http.StatusOK, invoiceView)
}
// CreateInvoice
// @Summary Create invoice
// @Description Create a new invoice
// @Tags invoices
// @Produce json
// @Param invoice body model.Invoice true "Invoice data"
// @Success 200 {object} model.InsertOneResult
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /invoices [post]
func CreateInvoice(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var invoice model.Invoice

	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var order model.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
	if err != nil {
		msg := fmt.Sprintf("message: Order was not found")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}
	status := "PENDING"
	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}

	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()

	validationErr := validate.Struct(invoice)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": validationErr.Error(),
		})
	}

	result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
	if insertErr != nil {
		msg := fmt.Sprintf("invoice item was not created")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}
// UpdateInvoice
// @Summary Update invoice
// @Description Update an invoice
// @Tags invoices
// @Produce json
// @Param invoice_id path string true "Invoice ID"
// @Param invoice body model.Invoice true "Invoice data"
// @Success 200 {object} model.UpdateResult
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /invoices/{invoice_id} [put]
func UpdateInvoice(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice model.Invoice
	invoiceId := c.Param("invoice_id")

	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	filter := bson.M{"invoice_id": invoiceId}

	var updateObj primitive.D

	if invoice.Payment_method != nil {
		updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
	}

	if invoice.Payment_status != nil {
		updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
	}

	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "PENDING"
	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}

	result, err := invoiceCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	if err != nil {
		msg := fmt.Sprintf("invoice item update failed")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}
