package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const acceptedRequestType = "application/json"

func RequestMustBeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestType := c.Request.Header.Get("Content-Type")
		if requestType != acceptedRequestType {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid content type",
				"code":    "ER-400",
			})
		}
	}
}
