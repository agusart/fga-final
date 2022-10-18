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
	photosController := controller.NewPhotosController(db)
	commentController := controller.NewCommentController(db)
	socialMediaController := controller.NewSocialMediaController(db)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	requestMustBeJSONMW := middleware.RequestMustBeJSON()
	authMW := middleware.Authentication()

	userRoute := router.Group("/users").Use(requestMustBeJSONMW)
	userRoute.POST("/register", userController.Register())
	userRoute.POST("/login", userController.Login())
	userRoute.Use(authMW).PUT("/", userController.Update())
	userRoute.DELETE("/", userController.Delete())

	photosRoute := router.Group("/photos").Use(authMW)
	photosRoute.GET("/", photosController.GetAll())
	photosRoute.DELETE("/:photoID", photosController.Delete())
	photosRoute.Use(requestMustBeJSONMW).POST("/", photosController.Create())
	photosRoute.PUT("/:photoID", photosController.Update())

	commentRoute := router.Group("/comments").Use(authMW)
	commentRoute.GET("/", commentController.Get())
	commentRoute.DELETE("/:commentID", commentController.Delete())
	commentRoute.Use(requestMustBeJSONMW).POST("/", commentController.Create())
	commentRoute.PUT("/:commentID", commentController.Update())

	socmedRoute := router.Group("/socialmedias").Use(authMW)
	socmedRoute.GET("/", socialMediaController.Get())
	socmedRoute.DELETE("/:socialMediaID", socialMediaController.Delete())
	socmedRoute.Use(requestMustBeJSONMW).POST("/", socialMediaController.Create())
	socmedRoute.PUT("/:socialMediaID", socialMediaController.Update())

	return router
}
