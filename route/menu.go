package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupMenuRoutes(e *echo.Group) {
	// e.Use(middlewares.AuthenticationMiddleware)
	menuRoutes := e.Group("/menus")
	{
		menuRoutes.GET("", controller.GetMenus)
		menuRoutes.GET("/:menu_id", controller.GetMenu)
		menuRoutes.POST("", controller.CreateMenu)
		menuRoutes.PATCH("/:menu_id", controller.UpdateMenu)
	}
}
