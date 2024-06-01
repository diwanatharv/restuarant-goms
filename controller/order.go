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

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
// GetOrders
// @Summary List orders
// @Description Get all orders
// @Tags orders
// @Produce json
// @Success 200 {array} model.CustomBsonM
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [get]
func GetOrders(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing order items",
		})
	}
	var allOrders []bson.M
	if err = result.All(ctx, &allOrders); err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, allOrders)
}
// GetOrder
// @Summary Get order
// @Description Get an order by ID
// @Tags orders
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{order_id} [get]
func GetOrder(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	orderId := c.Param("order_id")
	var order model.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while fetching the orders",
		})
	}
	return c.JSON(http.StatusOK, order)
}
// CreateOrder
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Produce json
// @Param order body model.Order true "Order data"
// @Success 200 {object} model.InsertOneResult
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [post]
func CreateOrder(c echo.Context) error {
	var table model.Table
	var order model.Order

	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validationErr := validate.Struct(order)

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": validationErr.Error(),
		})
	}

	if order.Table_id != nil {
		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("message:Table was not found")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": msg,
			})
		}
	}

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	result, insertErr := orderCollection.InsertOne(ctx, order)

	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	defer cancel()
	return c.JSON(http.StatusOK, result)
}
// UpdateOrder
// @Summary Update order
// @Description Update an order
// @Tags orders
// @Produce json
// @Param order_id path string true "Order ID"
// @Param order body model.Order true "Order data"
// @Success 200 {object} model.UpdateResult
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{order_id} [put]
func UpdateOrder(c echo.Context) error {
	var table model.Table
	var order model.Order

	var updateObj primitive.D

	orderId := c.Param("order_id")
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if order.Table_id != nil {
		err := menuCollection.FindOne(ctx, bson.M{"tabled_id": order.Table_id}).Decode(&table)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("message:Menu was not found")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": msg,
			})
		}
		updateObj = append(updateObj, bson.E{"menu", order.Table_id})
	}

	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

	upsert := true

	filter := bson.M{"order_id": orderId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		msg := fmt.Sprintf("order item update failed")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	defer cancel()
	return c.JSON(http.StatusOK, result)
}

func OrderItemOrderCreator(order model.Order) string {
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id
}
