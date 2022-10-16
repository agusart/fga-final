package controller

import (
	"fga-final/model"
	"fga-final/pkg/crypto"
	"fga-final/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type userController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *userController {
	return &userController{
		db: db,
	}
}

type registerResponse struct {
	Age      int    `json:"age"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (controller *userController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		err := controller.db.Create(&user).Error
		if err != nil {
			if strings.Contains(err.Error(), "idx_email") {
				err = errors.New("email already registered")
			}

			if strings.Contains(err.Error(), "idx_username") {
				err = errors.New("username already registered")
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		c.JSON(http.StatusCreated, registerResponse{
			Age:      user.Age,
			Email:    user.Email,
			ID:       user.ID,
			Username: user.Username,
		})
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (controller *userController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest loginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		var user model.User

		err := controller.db.Where("email = ?", loginRequest.Email).Take(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email / password",
				"error":   UnauthorizedError,
			})

			return
		}

		isPasswordvalid := crypto.CheckPasswordHash(loginRequest.Password, user.Password)
		if !isPasswordvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email / password",
				"error":   UnauthorizedError,
			})

			return
		}

		token, err := jwt.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"error":   UnauthorizedError,
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
