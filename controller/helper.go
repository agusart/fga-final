package controller

import (
	"fga-final/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) uint {
	id := c.MustGet(middleware.UserDataKey).(jwt.MapClaims)["id"].(float64)

	return uint(id)
}
