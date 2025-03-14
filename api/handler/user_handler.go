package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/utils"
	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
	"github.com/rmsubekti/indico/core/service"
	postgre "github.com/rmsubekti/indico/postgres"
)

type UserHandler struct {
	pg postgre.IPostgre
}

func NewUserHandler(pg postgre.IPostgre) UserHandler {
	return UserHandler{pg: pg}
}

func (u *UserHandler) Login(c *gin.Context) {
	var (
		ulogin port.UserLogin
		user   domain.User
		err    error
	)

	if c.ShouldBind(&ulogin) != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			"payload tidak sesuai",
		)
		return
	}

	tx, _ := u.pg.Begin()
	userSrv := service.NewUserService(tx.UserRepo())
	defer tx.Commit()

	if user, err = userSrv.Login(c, ulogin); err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"code":    http.StatusNotFound,
				"message": err.Error(),
			},
		)
		return
	}

	claim := utils.Claim{
		ID:         user.ID,
		Role:       string(user.UserRole),
		ExpireDays: 2,
	}

	if err = claim.CreateToken(); err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"code":    http.StatusUnauthorized,
				"message": "error membuat token",
			},
		)
		return
	}

	c.JSON(http.StatusOK, claim)
}

func (u *UserHandler) Register(c *gin.Context) {
	var (
		ureg port.UserRegister
		err  error
	)

	if c.ShouldBind(&ureg) != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}

	tx, _ := u.pg.Begin()
	userSrv := service.NewUserService(tx.UserRepo())
	defer tx.Commit()

	if err = userSrv.Register(c, &ureg); err != nil {
		tx.Rollback()
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"code":    http.StatusNotFound,
				"message": err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "user sukses didaftarkan",
		},
	)
}

func (u *UserHandler) GetMe(c *gin.Context) {
	claim := c.MustGet("claim").(utils.Claim)
	user := domain.User{ID: claim.ID}
	tx, _ := u.pg.Begin()
	userSrv := service.NewUserService(tx.UserRepo())
	defer tx.Commit()

	if err := userSrv.Get(c, &user); err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"code":    http.StatusNotFound,
				"message": http.StatusText(http.StatusNotFound),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"data": user,
		},
	)
}

func (u *UserHandler) List(c *gin.Context) {
	var users port.UserList
	tx, _ := u.pg.Begin()
	userSrv := service.NewUserService(tx.UserRepo())
	defer tx.Commit()

	if c.BindJSON(&users) != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}

	if err := userSrv.List(c, &users); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	c.JSON(http.StatusOK, users)
}
