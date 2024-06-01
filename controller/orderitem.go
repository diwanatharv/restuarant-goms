package controller

import (
	"context"
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

type OrderItemPack struct {
	Table_id    *string
	Order_items []model.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")
// GetOrderItems godoc
// @Summary Get all order items
// @Description Get all order items from the database
// @Tags orderItems
// @Accept  json
// @Produce  json
// @Success 200 {array} model.CustomBsonM
// @Failure 500 {object} map[string]interface{}
// @Router /orderitems [get]
func GetOrderItems(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderItemCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing ordered items",
		})
	}
	var allOrderItems []bson.M
	if err = result.All(ctx, &allOrderItems); err != nil {
		log.Fatal(err)
		return err
	}
	return c.JSON(http.StatusOK, allOrderItems)
}
// GetOrderItemsByOrder godoc
// @Summary Get order items by order ID
// @Description Get order items from the database by order ID
// @Tags orderItems
// @Accept  json
// @Produce  json
// @Param order_id path string true "Order ID"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object}  map[string]interface{}
// @Router /orderitems/order/{order_id} [get]
func GetOrderItemsByOrder(c echo.Context) error {
	orderId := c.Param("order_id")

	allOrderItems, err := ItemsByOrder(orderId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing order items by order ID",
		})
	}
	return c.JSON(http.StatusOK, allOrderItems)
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}}}

	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{

			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err

}
// GetOrderItem godoc
// @Summary Get a single order item
// @Description Get a single order item from the database by its ID
// @Tags orderItems
// @Accept  json
// @Produce  json
// @Param order_item_id path string true "Order Item ID"
// @Success 200 {object} model.OrderItem
// @Failure 500 {object}  map[string]interface{}
// @Router /orderitems/{order_item_id} [get]
func GetOrderItem(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("order_item_id")
	var orderItem model.OrderItem

	err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing ordered item",
		})
	}
	return c.JSON(http.StatusOK, orderItem)
}
// UpdateOrderItem godoc
// @Summary Update an order item
// @Description Update an order item in the database by its ID
// @Tags orderItems
// @Accept  json
// @Produce  json
// @Param order_item_id path string true "Order Item ID"
// @Param orderItem body model.OrderItem true "Order Item"
// @Success 200 {object} model.UpdateResult
// @Failure 500 {object}  map[string]interface{}
// @Router /orderitems/{order_item_id} [put]
func UpdateOrderItem(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItem model.OrderItem

	orderItemId := c.Param("order_item_id")

	filter := bson.M{"order_item_id": orderItemId}

	var updateObj primitive.D

	if orderItem.Unit_price != nil {
		updateObj = append(updateObj, bson.E{"unit_price", *&orderItem.Unit_price})
	}

	if orderItem.Quantity != nil {
		updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
	}

	if orderItem.Food_id != nil {
		updateObj = append(updateObj, bson.E{"food_id", *orderItem.Food_id})
	}

	orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		msg := "Order item update failed"
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}
// CreateOrderItem godoc
// @Summary Create new order items
// @Description Create new order items in the database
// @Tags orderItems
// @Accept  json
// @Produce  json
// @Param orderItemPack body OrderItemPack true "Order Item Pack"
// @Success 200 {object} model.InsertOneResult
// @Failure 400 {object} map[string]interface{}
// @Failure 500
func CreateOrderItem(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItemPack OrderItemPack
	var order model.Order

	if err := c.Bind(&orderItemPack); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	orderItemsToBeInserted := []interface{}{}
	order.Table_id = orderItemPack.Table_id
	order_id := OrderItemOrderCreator(order)

	for _, orderItem := range orderItemPack.Order_items {
		orderItem.Order_id = order_id

		validationErr := validate.Struct(orderItem)

		if validationErr != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": validationErr.Error(),
			})
		}

		orderItem.ID = primitive.NewObjectID()
		orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Order_item_id = orderItem.ID.Hex()
		var num = toFixed(*orderItem.Unit_price, 2)
		orderItem.Unit_price = &num
		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, insertedOrderItems)
}
