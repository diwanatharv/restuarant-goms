package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo) {
	// e.Use(middlewares.AuthenticationMiddleware)
	userRoutes := e.Group("/users")
	{
		userRoutes.GET("", controller.GetUsers)
		userRoutes.GET("/:user_id", controller.GetUser)
		userRoutes.POST("/signup", controller.SignUp)
		userRoutes.POST("/login", controller.Login)
	}
}
