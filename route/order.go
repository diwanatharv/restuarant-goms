package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupOrderRoutes(e *echo.Group) {
	// e.Use(middlewares.AuthenticationMiddleware)
	orderRoutes := e.Group("/orders")
	{
		orderRoutes.GET("", controller.GetOrders)
		orderRoutes.GET("/:order_id", controller.GetOrder)
		orderRoutes.POST("", controller.CreateOrder)
		orderRoutes.PATCH("/:order_id", controller.UpdateOrder)
	}
}
