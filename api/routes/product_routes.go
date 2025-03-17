package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/handler"
	"github.com/rmsubekti/indico/api/middleware"
	"github.com/rmsubekti/indico/core/domain"
	postgre "github.com/rmsubekti/indico/postgres"
)

func ProductRoutes(g *gin.Engine, pg postgre.IPostgre) {
	prodHandler := handler.NewProductHandler(pg)
	products := g.Group("/products")
	products.Use(middleware.Auth())
	products.POST("", middleware.Role(domain.UserAdmin), prodHandler.Create)
	products.GET("", prodHandler.List)
	products.GET("/:id", prodHandler.Get)
	products.PUT("/:id", middleware.Role(domain.UserAdmin), prodHandler.Update)
	products.DELETE("/:id", middleware.Role(domain.UserAdmin), prodHandler.Delete)
}
