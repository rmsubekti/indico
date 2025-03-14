package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/routes"
	"github.com/rmsubekti/indico/config"
	postgre "github.com/rmsubekti/indico/postgres"
)

func main() {
	pg := postgre.NewPostgre(config.PG_DSN())
	defer pg.Close()
	if err := pg.Migrate(); err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	routes.UserRoutes(r, pg)
	r.Run(":8080")

}
