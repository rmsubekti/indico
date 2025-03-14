package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/handler"
	"github.com/rmsubekti/indico/api/middleware"
	"github.com/rmsubekti/indico/core/domain"
	postgre "github.com/rmsubekti/indico/postgres"
)

func UserRoutes(g *gin.Engine, pg postgre.IPostgre) {
	userHandler := handler.NewUserHandler(pg)
	g.POST("/login", userHandler.Login)
	g.POST("/register", userHandler.Register)

	user := g.Group("/users")
	user.Use(middleware.Auth())

	user.GET("/me", userHandler.GetMe)
	user.GET("", middleware.Role(domain.UserAdmin), userHandler.List)
}
