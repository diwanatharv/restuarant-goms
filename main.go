package main

import (
	"management/database"
	middlewares "management/middleware"
	"management/route"
	"os"

	_ "management/docs" // You need to import the generated docs package

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// @title Restaurant Management API
// @version 1.0
// @description This is a sample server for a restaurant management system.
// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// Create a new Echo instance
	e := echo.New()

	// Use the logger middleware
	e.Use(middleware.Logger())

	// Set up routes
	route.SetupUserRoutes(e)
	authRoutes := e.Group("")
	authRoutes.Use(middlewares.AuthenticationMiddleware)
	route.SetupFoodRoutes(authRoutes)
	route.SetupMenuRoutes(authRoutes)
	route.SetupTableRoutes(authRoutes)
	route.SetupOrderRoutes(authRoutes)
	route.SetupOrderItemRoutes(authRoutes)
	route.SetupInvoiceRoutes(authRoutes)

	// Redirect the base URL to Swagger UI
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(301, "/swagger/index.html")
	})

	// Start the server
	e.Logger.Fatal(e.Start(":" + port))
}
