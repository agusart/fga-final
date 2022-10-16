package jwt

import "github.com/dgrijalva/jwt-go"

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
