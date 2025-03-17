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

// Login godoc
// @Summary Login
// @Description login user
// @Tags         user
// @Produce  json
// @Param request body string true " Body payload message/rfc822" SchemaExample({\n\t"email": "email",\n\t"password": "passwrod"\n})
// @Success 200 {object} utils.Claim
// @Router /login [post]
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

// Register godoc
// @Summary Register
// @Description register user
// @Tags         user
// @Produce  json
// @Param request body string true " Body payload message/rfc822" SchemaExample({\n\t"name": "name",\n\t"email": "email",\n\t"password": "passwrod",\n\t"confirm_password": "passwrod",\n\t"user_role": "Admin|Staff"\n})
// @Success 200 {string} success
// @Router /register [post]
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

// GetMe godoc
// @Summary Get login info
// @Tags         user
// @Produce  json
// @Success 200 {object} domain.User
// @Router /users/me [get]
// @Security Bearer
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

// List godoc
// @Summary Get user list
// @Tags         user
// @Produce  json
// @Success 200 {object} port.UserList
// @Param        page    query    int  false  "show data on page n"
// @Param        limit    query     int  false  "limit items per page"
// @Param        search    query     string  false  "search  filter by name"
// @Router /users [get]
// @Security Bearer
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
