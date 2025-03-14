package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
	"github.com/rmsubekti/indico/core/service"
	postgre "github.com/rmsubekti/indico/postgres"
)

type WarehouseHandler struct {
	pg postgre.IPostgre
}

func NewWarehouseHandler(pg postgre.IPostgre) WarehouseHandler {
	return WarehouseHandler{pg}
}

func (w *WarehouseHandler) Create(c *gin.Context) {
	var (
		warehouse domain.Warehouse
		err       error
	)

	c.ShouldBind(&warehouse)
	tx, _ := w.pg.Begin()
	waSrv := service.NewWarehouseService(tx.WarehouseRepo())
	defer tx.Commit()

	if err = waSrv.Add(c, &warehouse); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

func (w *WarehouseHandler) List(c *gin.Context) {
	var (
		list port.WarehouseList
		err  error
	)
	c.ShouldBind(&list)
	tx, _ := w.pg.Begin()
	waSrv := service.NewWarehouseService(tx.WarehouseRepo())
	defer tx.Commit()

	if err = waSrv.List(c, &list); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, list)
}
