package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
	"github.com/rmsubekti/indico/core/service"
	postgre "github.com/rmsubekti/indico/postgres"
)

type OrderHandler struct {
	pg postgre.IPostgre
}

func NewOrderHandler(pg postgre.IPostgre) OrderHandler {
	return OrderHandler{pg: pg}
}

func (o *OrderHandler) Shipping(c *gin.Context) {
	var order domain.Order
	c.ShouldBind(&order)
	order.OrderType = domain.TypeShipping
	if err := createOrder(c, order, o.pg); err != nil {
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}
	c.JSON(http.StatusOK, "shipping diproses")
}
func (o *OrderHandler) Receive(c *gin.Context) {
	var order domain.Order
	c.ShouldBind(&order)
	order.OrderType = domain.TypeReceive
	if err := createOrder(c, order, o.pg); err != nil {
		c.AbortWithStatusJSON(http.StatusNotModified, err)
		return
	}
	c.JSON(http.StatusOK, "receive diproses")
}

func createOrder(ctx context.Context, order domain.Order, pg postgre.IPostgre) (err error) {

	tx, _ := pg.Begin()
	detServ := service.NewOrderDetailService(tx.OrderDetailRepo())
	orderSrv := service.NewOrderService(tx.OrderRepo()).WithOrderDetailService(detServ)

	if err = orderSrv.Add(ctx, &order); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	//simulasi
	var delaySecond uint = 15
	go func(pg postgre.IPostgre, order_id uint, delay uint) {
		time.Sleep(time.Second * time.Duration(delay))
		simulasi(pg, order_id, domain.StatusOrderProcessing)
		time.Sleep(time.Second * time.Duration(delay))
		simulasi(pg, order_id, domain.StatusOrderCompleted)
	}(pg, order.ID, delaySecond)

	return
}

func simulasi(pg postgre.IPostgre, order_id uint, status domain.OrderStatus) {
	tx, _ := pg.Begin()
	var err error
	ctx := context.Background()
	defer ctx.Done()
	detServ := service.NewOrderDetailService(tx.OrderDetailRepo())
	stockSrv := service.NewStockService(tx.StockRepo())
	orderSrv := service.NewOrderService(tx.OrderRepo()).WithOrderDetailService(detServ).WithStockService(stockSrv)
	if err = orderSrv.ChangeStatus(ctx, order_id, status); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}

func (o *OrderHandler) List(c *gin.Context) {
	var orderList port.OrderList
	c.ShouldBind(&orderList)
	tx, _ := o.pg.Begin()
	whSrv := service.NewWarehouseService(tx.WarehouseRepo())
	orderSrv := service.NewOrderService(tx.OrderRepo()).WithWarehouseService(whSrv)

	if err := orderSrv.List(c, &orderList); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, orderList)
}

func (o *OrderHandler) GetDetailed(c *gin.Context) {
	var (
		order domain.Order
		id, _ = strconv.Atoi(c.Param("id"))
		err   error
	)
	tx, _ := o.pg.Begin()
	whSrv := service.NewWarehouseService(tx.WarehouseRepo())
	detServ := service.NewOrderDetailService(tx.OrderDetailRepo())
	prodServ := service.NewProductService(tx.ProductRepo())
	orderSrv := service.NewOrderService(tx.OrderRepo()).WithWarehouseService(whSrv).WithOrderDetailService(detServ).WithProductService(prodServ)

	if order, err = orderSrv.GetByID(c, uint(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, order)
}
