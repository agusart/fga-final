package controller

import (
	"fga-final/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type socialMediaController struct {
	db *gorm.DB
}

func NewSocialMediaController(db *gorm.DB) *socialMediaController {
	return &socialMediaController{
		db: db,
	}
}

func (controller *socialMediaController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var socialMedia model.SocialMedia

		photoID := c.Param("socialMediaID")
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
			First(&socialMedia).RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})

			return
		}

		if socialMedia.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   UnauthorizedError,
			})

			return
		}

		if err := controller.db.Delete(&socialMedia).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your social media has been successfully deleted",
		})

	}
}

type createSocialMediaResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         uint   `json:"user_id"`
	CreatedAt      string `json:"created_at"`
}

func (controller *socialMediaController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var socialMedia model.SocialMedia
		if err := c.ShouldBindJSON(&socialMedia); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		socialMedia.UserID = GetUserIDFromContext(c)
		err := controller.db.Create(&socialMedia).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusCreated, createSocialMediaResponse{
			ID:             socialMedia.ID,
			Name:           socialMedia.Name,
			SocialMediaURL: socialMedia.SocialMediaURL,
			UserID:         socialMedia.UserID,
			CreatedAt:      socialMedia.CreatedAt.Format(time.RFC3339),
		})
	}
}

type updateSocialMediaResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         uint   `json:"user_id"`
	UpdatedAt      string `json:"updated_at"`
}

func (controller *socialMediaController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var socialMedia model.SocialMedia

		socialMediaID := c.Param("socialMediaID")
		if socialMediaID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
				"error":   BadRequestError,
			})

			return
		}

		if controller.db.
			Preload("User").
			Where("id = ?", socialMediaID).
			First(&socialMedia).RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})

			return
		}

		if socialMedia.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
				"error":   UnauthorizedError,
			})

			return
		}

		if err := c.ShouldBindJSON(&socialMedia); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		err := controller.db.Save(&socialMedia).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, updateSocialMediaResponse{
			ID:             socialMedia.ID,
			Name:           socialMedia.Name,
			SocialMediaURL: socialMedia.SocialMediaURL,
			UserID:         socialMedia.UserID,
			UpdatedAt:      socialMedia.CreatedAt.Format(time.RFC3339),
		})

	}
}

type getSocialMediaResponse struct {
	ID             uint                       `json:"id"`
	Name           string                     `json:"name"`
	SocialMediaURL string                     `json:"social_media_url"`
	UserID         uint                       `json:"UserId"`
	CreatedAt      string                     `json:"createdAt"`
	UpdatedAt      string                     `json:"updatedAt"`
	User           getSocialMediaUserResponse `json:"User"`
}

type getSocialMediaUserResponse struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (controller *socialMediaController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			socialMedias []model.SocialMedia
		)

		socialMediaResponse := make([]getSocialMediaResponse, 0)

		if err := controller.db.Preload("User").Where("user_id = ?", GetUserIDFromContext(c)).Find(&socialMedias).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   InternalServerError,
			})

			return
		}

		for _, socialMedia := range socialMedias {
			socialMediaResponse = append(socialMediaResponse, getSocialMediaResponse{
				ID:             socialMedia.ID,
				Name:           socialMedia.Name,
				SocialMediaURL: socialMedia.SocialMediaURL,
				UserID:         socialMedia.UserID,
				CreatedAt:      socialMedia.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      socialMedia.UpdatedAt.Format(time.RFC3339),
				User: getSocialMediaUserResponse{
					Username: socialMedia.User.Username,
					ID:       socialMedia.User.ID,
				},
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"social_medias": socialMediaResponse,
		})

	}
}
