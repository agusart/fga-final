package controller

import (
	"fga-final/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type commentController struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *commentController {
	return &commentController{
		db: db,
	}
}

type createCommentResponse struct {
	ID        uint   `json:"id"`
	Message   string `json:"message"`
	PhotoID   uint   `json:"photo_id"`
	UserID    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (controller *commentController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment model.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		comment.UserID = GetUserIDFromContext(c)
		err := controller.db.Create(&comment).Error
		if err != nil {
			if strings.Contains(err.Error(), "fk_comments_photo") {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "photo not found",
					"error":   NotFoundError,
				})

				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, createCommentResponse{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		})
	}
}

type updateCommentResponse struct {
	ID        uint   `json:"id"`
	Message   string `json:"message"`
	PhotoID   uint   `json:"photo_id"`
	UserID    uint   `json:"user_id"`
	UpdatedAt string `json:"upadated_at"`
}

func (controller *commentController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment model.Comment

		commentID := c.Param("commentID")
		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
				"error":   BadRequestError,
			})

			return
		}

		controller.db.
			Preload("User").
			Where("id = ?", commentID).
			First(&comment)

		if comment.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})

			return
		}

		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		err := controller.db.Save(&comment).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, updateCommentResponse{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		})

	}
}

func (controller *commentController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment model.Comment

		photoID := c.Param("commentID")
		if photoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
				"error":   BadRequestError,
			})

			return
		}

		controller.db.
			Preload("User").
			Where("id = ?", photoID).
			First(&comment)

		if comment.UserID != GetUserIDFromContext(c) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"error":   NotFoundError,
			})

			return
		}

		if err := controller.db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your comment has been successfully deleted",
		})

	}
}

type getCommentResponse struct {
	ID        uint                    `json:"id"`
	Message   string                  `json:"message"`
	PhotoID   uint                    `json:"photo_id"`
	UserID    uint                    `json:"user_id"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt string                  `json:"updated_at"`
	User      getCommentUserResponse  `json:"User"`
	Photo     getCommentPhotoResponse `json:"Photo"`
}

type getCommentUserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"message"`
	Username string `json:"username"`
}

type getCommentPhotoResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
}

func (controller *commentController) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comments []model.Comment
		commentResponses := make([]getCommentResponse, 0)

		err := controller.db.
			Joins("inner join users on users.id = comments.user_id and users.deleted_at is null").
			Joins("inner join photos on photos.id = comments.photo_id and photos.deleted_at is null").
			Preload("User").
			Preload("Photo").
			Find(&comments).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		for _, comment := range comments {
			commentResponse := getCommentResponse{
				ID:        comment.ID,
				Message:   comment.Message,
				PhotoID:   comment.PhotoID,
				UserID:    comment.UserID,
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
				UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
			}

			commentPhoto := getCommentPhotoResponse{
				ID:       comment.Photo.ID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoURL: comment.Photo.PhotoURL,
				UserID:   comment.Photo.UserID,
			}

			commentUser := getCommentUserResponse{
				Email:    comment.User.Email,
				Username: comment.User.Username,
				ID:       comment.User.ID,
			}

			commentResponse.Photo = commentPhoto
			commentResponse.User = commentUser

			commentResponses = append(commentResponses, commentResponse)
		}

		c.JSON(http.StatusOK, commentResponses)
	}
}
