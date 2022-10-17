package controller

import (
	"fga-final/middleware"
	"fga-final/pkg/jwthelper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) uint {
	userID := c.MustGet(middleware.UserDataKey).(jwt.MapClaims)[jwthelper.IdKey].(uint)

	return userID
}
