package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmsubekti/indico/api/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := utils.Claim{Token: c.GetHeader("Authorization")}
		if err := claim.Parse(); err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"code":    http.StatusUnauthorized,
					"message": err.Error(),
				},
			)
			return
		}

		c.Set("claim", claim)
		c.Next()
	}
}
