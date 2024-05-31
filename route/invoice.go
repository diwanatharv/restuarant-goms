package route

import (
	"management/controller"

	"github.com/labstack/echo/v4"
)

func SetupInvoiceRoutes(g *echo.Group) {
	// e.Use(middlewares.AuthenticationMiddleware)
	invoiceRoutes := g.Group("/invoices")

	invoiceRoutes.GET("", controller.GetInvoices)
	invoiceRoutes.GET("/:invoice_id", controller.GetInvoice)
	invoiceRoutes.POST("", controller.CreateInvoice)
	invoiceRoutes.PATCH("/:invoice_id", controller.UpdateInvoice)

}
