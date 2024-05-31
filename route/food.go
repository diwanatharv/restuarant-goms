package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupFoodRoutes(g *echo.Group) {
    foodRoutes := g.Group("/foods")

    foodRoutes.GET("", controller.GetFoods)
    foodRoutes.GET("/:food_id", controller.GetFood)
    foodRoutes.POST("", controller.CreateFood)
    foodRoutes.PATCH("/:food_id", controller.UpdateFood)
}
