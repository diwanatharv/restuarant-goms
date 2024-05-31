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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := tableCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing table items",
		})
	}

	var allTables []bson.M
	if err = result.All(ctx, &allTables); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, allTables)
}

func GetTable(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	tableId := c.Param("table_id")
	var table model.Table

	err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while fetching the tables",
		})
	}

	return c.JSON(http.StatusOK, table)
}

func CreateTable(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table model.Table

	if err := c.Bind(&table); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	validationErr := validate.Struct(table)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": validationErr.Error(),
		})
	}

	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
	table.Table_id = table.ID.Hex()

	result, insertErr := tableCollection.InsertOne(ctx, table)
	if insertErr != nil {
		msg := fmt.Sprintf("Table item was not created")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateTable(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table model.Table
	tableId := c.Param("table_id")

	if err := c.Bind(&table); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var updateObj primitive.D

	if table.Number_of_guests != nil {
		updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
	}

	if table.Table_number != nil {
		updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
	}

	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	filter := bson.M{"table_id": tableId}
	result, err := tableCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		msg := fmt.Sprintf("table item update failed")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}
