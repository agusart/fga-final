package controller

import (
	"fga-final/middleware"
	"fga-final/model"
	"fga-final/pkg/crypto"
	"fga-final/pkg/jwthelper"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
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
	Email    string `json:"email" valid:"required~email is blank,email~invalid email format"`
	Password string `json:"password"  valid:"required~password is blank"`
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

		_, err := govalidator.ValidateStruct(loginRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		var user model.User

		err = controller.db.Where("email = ?", loginRequest.Email).Take(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email / password",
				"error":   UnauthorizedError,
			})

			return
		}

		isPasswordValid := crypto.CheckPasswordHash(loginRequest.Password, user.Password)
		if !isPasswordValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email / password",
				"error":   UnauthorizedError,
			})

			return
		}

		token, err := jwthelper.GenerateToken(user.ID)
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

type updateUserRequest struct {
	Email    string `json:"email" valid:"required~email is blank,email~invalid email format"`
	Username string `json:"username"  valid:"required~username is blank"`
}

type editResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Age       int       `json:"age"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (controller *userController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateUserRequest updateUserRequest
		if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json format",
				"error":   BadRequestError,
			})

			return
		}

		_, err := govalidator.ValidateStruct(updateUserRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   BadRequestError,
			})

			return
		}

		var user model.User
		id := c.MustGet(middleware.UserDataKey).(jwt.MapClaims)["id"]

		err = controller.db.Where("id = ?", id).Take(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   InternalServerError,
			})

			return
		}

		user.Username = updateUserRequest.Username
		user.Email = updateUserRequest.Email
		user.UpdatedAt = time.Now()

		if err := controller.db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   InternalServerError,
			})

			return
		}

		c.JSON(http.StatusOK, editResponse{
			ID:        user.ID,
			Email:     user.Email,
			Age:       user.Age,
			Username:  user.Username,
			UpdatedAt: user.UpdatedAt,
		})

	}
}
