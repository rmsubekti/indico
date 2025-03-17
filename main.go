package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/routes"
	"github.com/rmsubekti/indico/config"
	"github.com/rmsubekti/indico/docs"
	postgre "github.com/rmsubekti/indico/postgres"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           INDICO Service
// @version         1.0
// @description     Login to create token.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	pg := postgre.NewPostgre(config.PG_DSN())
	defer pg.Close()
	if err := pg.Migrate(); err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	cfg := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(cfg))
	docs.SwaggerInfo.BasePath = "/"

	routes.UserRoutes(r, pg)
	routes.ProductRoutes(r, pg)
	routes.LocationRoutes(r, pg)
	routes.OrderRoutes(r, pg)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":" + config.APP.PORT)

}
