package jwthelper

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

const secret = "secret"

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id": userID,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := parseToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	errResp := errors.New("Sign To Proceed")
	headerToken := c.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errResp
	}

	stringToken := strings.Split(headerToken, " ")[1]
	token, _ := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResp
		}

		return []byte(secret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResp
	}

	return token.Claims.(jwt.MapClaims), nil
}
