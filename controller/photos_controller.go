package controller

import (
	"fga-final/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type photosController struct {
	db *gorm.DB
}

func NewPhotosController(db *gorm.DB) *photosController {
	return &photosController{
		db: db,
	}
}

type createPhotoResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	PhotoURL  string `json:"photo_url"`
	UserID    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (controller *photosController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var photo model.Photo
		if err := c.ShouldBindJSON(&photo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		userID := GetUserIDFromContext(c)
		photo.UserID = userID

		err := controller.db.Create(&photo).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusCreated, createPhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt.Format(time.RFC3339),
		})
	}
}

type getPhotoResponse struct {
	ID        uint                 `json:"id"`
	Title     string               `json:"title"`
	Caption   string               `json:"caption"`
	PhotoURL  string               `json:"photo_url"`
	UserID    uint                 `json:"user_id"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
	User      getPhotoUserResponse `json:"User"`
}

type getPhotoUserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (controller *photosController) GetMine() gin.HandlerFunc {
	return func(c *gin.Context) {
		var photos []model.Photo

		err := controller.db.Where("user_id = ?", GetUserIDFromContext(c)).Find(&photos).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		photosResponse := make([]getPhotoResponse, 0)
		for _, photo := range photos {
			photosResponse = append(photosResponse, getPhotoResponse{
				ID:        photo.ID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				PhotoURL:  photo.PhotoURL,
				UserID:    photo.UserID,
				CreatedAt: photo.CreatedAt.Format(time.RFC3339),
				UpdatedAt: photo.UpdatedAt.Format(time.RFC3339),
				User: getPhotoUserResponse{
					Email:    photo.User.Email,
					Username: photo.User.Username,
				},
			})
		}

		c.JSON(http.StatusOK, photosResponse)
	}
}

func (controller *photosController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var photos []model.Photo

		err := controller.db.
			Joins("inner join users on users.id = photos.user_id and users.deleted_at is null").
			Preload("User").Find(&photos).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		photosResponse := make([]getPhotoResponse, 0)
		for _, photo := range photos {
			photosResponse = append(photosResponse, getPhotoResponse{
				ID:        photo.ID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				PhotoURL:  photo.PhotoURL,
				UserID:    photo.UserID,
				CreatedAt: photo.CreatedAt.Format(time.RFC3339),
				UpdatedAt: photo.UpdatedAt.Format(time.RFC3339),
				User: getPhotoUserResponse{
					Email:    photo.User.Email,
					Username: photo.User.Username,
				},
			})
		}

		c.JSON(http.StatusOK, photosResponse)
	}
}

type updatePhotoResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	PhotoURL  string `json:"photo_url"`
	UserID    uint   `json:"user_id"`
	UpdatedAt string `json:"updated_at"`
}

func (controller *photosController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var photo model.Photo

		photoID := c.Param("photoID")
		if photoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
				"error":   BadRequestError,
			})

			return
		}

		if controller.db.
			Preload("User").
			Where("id = ?", photoID).
			First(&photo).RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})

			return
		}

		if photo.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   UnauthorizedError,
			})

			return
		}

		if err := c.ShouldBindJSON(&photo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		err := controller.db.Save(&photo).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, updatePhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    photo.UserID,
			UpdatedAt: photo.UpdatedAt.Format(time.RFC3339),
		})

	}
}

func (controller *photosController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var photo model.Photo

		photoID := c.Param("photoID")
		if photoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
				"error":   BadRequestError,
			})

			return
		}

		if controller.db.
			Preload("User").
			Where("id = ?", photoID).
			First(&photo).RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})
			return
		}

		if photo.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   UnauthorizedError,
			})
			return
		}

		if err := controller.db.Delete(&photo).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your photo has been successfully deleted",
		})

	}
}
