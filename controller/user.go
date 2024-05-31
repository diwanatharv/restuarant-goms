package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"management/database"
	"management/helper"
	"management/model"
	"net/http"
	"strconv"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, projectStage,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing user items",
		})
	}

	var allUsers []bson.M
	if err = result.All(ctx, &allUsers); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while decoding user items",
		})
	}

	return c.JSON(http.StatusOK, allUsers[0])
}

func GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userId := c.Param("user_id")
	var user model.User

	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error occurred while listing user items",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func SignUp(c echo.Context) error {
	fmt.Println("under user sign up")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": validationErr.Error()})
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error occurred while checking for the email"})
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error occurred while checking for the phone number"})
	}

	if count > 0 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "this email or phone number already exists"})
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
	user.Token = &token
	user.Refresh_Token = &refreshToken

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		msg := fmt.Sprintf("User item was not created")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": msg})
	}

	return c.JSON(http.StatusOK, result)
}

func Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var foundUser model.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "user not found, login seems to be incorrect"})
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	if !passwordIsValid {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": msg})
	}

	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
	helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	return c.JSON(http.StatusOK, foundUser)
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	return check, msg
}
