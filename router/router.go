package router

import (
	"fga-final/controller"
	"fga-final/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func StartRouter(db *gorm.DB) *gin.Engine {
	userController := controller.NewUserController(db)
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	router.Use(middleware.RequestMustBeJSON()).POST("/users/register", userController.Register())
	router.Use(middleware.RequestMustBeJSON()).POST("/users/login", userController.Login())

	return router
}
