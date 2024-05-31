package main

import (
	"management/database"
	middlewares "management/middleware"
	"management/route"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// Create a new Echo instance
	e := echo.New()

	// // Use the logger middleware
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

	// Start the server
	e.Logger.Fatal(e.Start(":" + port))
}
