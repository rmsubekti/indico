package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/handler"
	"github.com/rmsubekti/indico/api/middleware"
	"github.com/rmsubekti/indico/core/domain"
	postgre "github.com/rmsubekti/indico/postgres"
)

func LocationRoutes(g *gin.Engine, pg postgre.IPostgre) {
	whHandler := handler.NewWarehouseHandler(pg)
	wh := g.Group("/locations")
	wh.Use(middleware.Auth())
	wh.POST("", middleware.Role(domain.UserAdmin), whHandler.Create)
	wh.GET("", whHandler.List)
}
