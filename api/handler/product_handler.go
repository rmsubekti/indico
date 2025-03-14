package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
	"github.com/rmsubekti/indico/core/service"
	postgre "github.com/rmsubekti/indico/postgres"
)

type ProductHandler struct {
	pg postgre.IPostgre
}

func NewProductHandler(pg postgre.IPostgre) ProductHandler {
	return ProductHandler{pg}
}

func (p *ProductHandler) Create(c *gin.Context) {
	var (
		product domain.Product
		err     error
	)
	c.ShouldBind(&product)

	tx, _ := p.pg.Begin()
	productSrv := service.NewProductService(tx.ProductRepo())
	defer tx.Commit()

	if err = productSrv.Add(c, &product); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (p *ProductHandler) Update(c *gin.Context) {
	var (
		update port.ProductUpdate
		err    error
		id, _  = strconv.Atoi(c.Param("id"))
	)
	c.ShouldBind(&update)

	update.ID = uint(id)

	tx, _ := p.pg.Begin()
	productSrv := service.NewProductService(tx.ProductRepo())
	defer tx.Commit()

	if err = productSrv.Update(c, &update); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, "product terupdate")
}

func (p *ProductHandler) List(c *gin.Context) {
	var (
		list port.ProductList
		err  error
	)
	c.ShouldBind(&list)

	tx, _ := p.pg.Begin()
	productSrv := service.NewProductService(tx.ProductRepo())
	defer tx.Commit()

	if err = productSrv.List(c, &list); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (p *ProductHandler) Get(c *gin.Context) {
	var (
		product domain.Product
		err     error
		id, _   = strconv.Atoi(c.Param("id"))
	)
	tx, _ := p.pg.Begin()
	productSrv := service.NewProductService(tx.ProductRepo())
	defer tx.Commit()

	if product, err = productSrv.Get(c, uint(id)); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (p *ProductHandler) Delete(c *gin.Context) {
	var (
		err   error
		id, _ = strconv.Atoi(c.Param("id"))
	)
	tx, _ := p.pg.Begin()
	productSrv := service.NewProductService(tx.ProductRepo())
	defer tx.Commit()

	if err = productSrv.Delete(c, uint(id)); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}

	c.JSON(http.StatusOK, "produk dihapus")
}
