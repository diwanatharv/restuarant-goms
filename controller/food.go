package controller

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"management/database"
	"management/model"
	"math"
	"net/http"
	"strconv"
	"time"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordPerPage, err := strconv.Atoi(c.QueryParam("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.QueryParam("startIndex"))

	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
	projectStage := bson.D{
		{
			"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

	result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occured while listing food items",
		})
	}

	var allFoods []bson.M
	if err = result.All(ctx, &allFoods); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, allFoods[0])
}
func GetFood(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	foodId := c.Param("food_id")
	var food model.Food
	err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
	defer cancel()
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, "error:occured while fewtching the food item")
	}
	return c.JSON(http.StatusOK, food)
}
func CreateFood(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var menu model.Menu
	var food model.Food

	if err := c.Bind(&food); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	validationErr := validate.Struct(food)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, validationErr.Error())

	}
	err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
	defer cancel()
	if err != nil {
		msg := fmt.Sprintf("menu was not found")
		return c.JSON(http.StatusInternalServerError, msg)

	}
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.Food_id = food.ID.Hex()
	var num = toFixed(*food.Price, 2)
	food.Price = &num

	result, insertErr := foodCollection.InsertOne(ctx, food)
	if insertErr != nil {
		msg := fmt.Sprintf("Food item was not created")
		return c.JSON(http.StatusInternalServerError, msg)

	}
	defer cancel()
	return c.JSON(http.StatusOK, result)
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func UpdateFood(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu model.Menu
	var food model.Food

	foodId := c.Param("food_id")

	if err := c.Bind(&food); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var updateObj primitive.D

	if food.Name != nil {
		updateObj = append(updateObj, bson.E{"name", food.Name})
	}

	if food.Price != nil {
		updateObj = append(updateObj, bson.E{"price", food.Price})
	}

	if food.Food_image != nil {
		updateObj = append(updateObj, bson.E{"food_image", food.Food_image})
	}

	if food.Menu_id != nil {
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			msg := fmt.Sprintf("message:Menu was not found")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": msg,
			})
		}
		updateObj = append(updateObj, bson.E{"menu", food.Price})
	}

	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

	upsert := true
	filter := bson.M{"food_id": foodId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := foodCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		msg := fmt.Sprint("food item update failed")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}

	return c.JSON(http.StatusOK, result)
}
