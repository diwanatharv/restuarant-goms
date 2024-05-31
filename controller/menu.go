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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := menuCollection.Find(context.TODO(), bson.M{})
	defer cancel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error occured while listing the menu items")
	}
	var allMenus []model.Menu
	if err = result.All(ctx, &allMenus); err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, allMenus)
}

func GetMenu(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	menuId := c.Param("menu_id")
	var menu model.Menu

	err := foodCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
	defer cancel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error occured while fetching the menu")
	}
	return c.JSON(http.StatusOK, menu)
}
func CreateMenu(c echo.Context) error {
	var menu model.Menu
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())

	}
	validationErr := validate.Struct(menu)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, validationErr.Error())
	}

	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.ID = primitive.NewObjectID()
	menu.Menu_id = menu.ID.Hex()

	result, insertErr := menuCollection.InsertOne(ctx, menu)
	if insertErr != nil {
		msg := fmt.Sprintf("Menu item was not created")
		return c.JSON(http.StatusInternalServerError, msg)

	}

	return c.JSON(http.StatusOK, result)
}
func UpdateMenu(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	menuId := c.Param("menu_id")
	var menu model.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())

	}
	validationErr := validate.Struct(menu)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, validationErr.Error())
	}
	filter := bson.M{"menu_id": menuId}
	update := bson.M{"$set": menu}
	opts := options.FindOneAndUpdate().SetUpsert(true)

	res, err := foodCollection.FindOneAndUpdate(ctx, filter, update, opts).Raw()
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	defer cancel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error occured while updating  the menu")
	}
	return c.JSON(http.StatusOK, res)
}
