package main

import (
	"log"
	"management/database"
	middlewares "management/middleware"
	"management/route"
	"net/http"
	"os"

	_ "management/docs" // You need to import the generated docs package

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
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

	// Serve the Swagger UI at /swagger/*
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

	// Cron job trigger endpoint
	e.GET("/run-cron-job", func(c echo.Context) error {
		// Trigger base URL logic here
		resp, err := http.Get("http://localhost:" + port + "/")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error triggering base URL: "+err.Error())
		}
		defer resp.Body.Close()
		return c.String(http.StatusOK, "Cron job executed successfully! Status: "+resp.Status)
	})

	// Set up and start the cron job
	c := cron.New()
	_, err := c.AddFunc("*/7 * * * *", func() {
		resp, err := http.Get("http://localhost:" + port + "/run-cron-job")
		if err != nil {
			log.Printf("Error triggering cron job: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Cron job executed successfully! Status: %s", resp.Status)
	})
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}
	c.Start()

	// Start the server
	e.Logger.Fatal(e.Start(":" + port))
}
