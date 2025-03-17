package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/handler"
	"github.com/rmsubekti/indico/api/middleware"
	"github.com/rmsubekti/indico/core/domain"
	postgre "github.com/rmsubekti/indico/postgres"
)

func OrderRoutes(g *gin.Engine, pg postgre.IPostgre) {
	orderHandler := handler.NewOrderHandler(pg)
	orders := g.Group("/orders")
	orders.Use(middleware.Auth())
	orders.POST("/receive", middleware.Role(domain.UserStaff), orderHandler.Receive)
	orders.POST("/ship", middleware.Role(domain.UserStaff), orderHandler.Shipping)
	orders.GET("", orderHandler.List)
	orders.GET("/:id", orderHandler.GetDetailed)
}
