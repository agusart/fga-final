package middleware

import (
	"fga-final/pkg/jwthelper"
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserDataKey = "userData"

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := jwthelper.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		}

		c.Set("userData", verifyToken)
		c.Next()
	}
}
