package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupOrderItemRoutes(e *echo.Group) {
	// e.Use(middlewares.AuthenticationMiddleware)
	orderItemRoutes := e.Group("/orderItems")
	{
		orderItemRoutes.GET("", controller.GetOrderItems)
		orderItemRoutes.GET("/:orderItem_id", controller.GetOrderItem)
		orderItemRoutes.GET("/order/:order_id", controller.GetOrderItemsByOrder)
		orderItemRoutes.POST("", controller.CreateOrderItem)
		orderItemRoutes.PATCH("/:orderItem_id", controller.UpdateOrderItem)
	}
}
