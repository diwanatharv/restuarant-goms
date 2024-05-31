package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupTableRoutes(e *echo.Group) {
	// e.Use(middlewares.AuthenticationMiddleware)
	tableRoutes := e.Group("/tables")
	{
		tableRoutes.GET("", controller.GetTables)
		tableRoutes.GET("/:table_id", controller.GetTable)
		tableRoutes.POST("", controller.CreateTable)
		tableRoutes.PATCH("/:table_id", controller.UpdateTable)
	}
}
