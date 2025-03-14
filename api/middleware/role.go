package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/utils"
	"github.com/rmsubekti/indico/core/domain"
)

func Role(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := c.MustGet("claim").(utils.Claim)
		if claim.Role != string(role) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"code":    http.StatusUnauthorized,
					"message": "role tidak sesuai",
				},
			)
			return
		}
		c.Next()
	}
}
