package helpers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenToUid(c *gin.Context) (int, error, string) {
	tokenString := c.Request.Header.Get("Authorization")
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 403, err, ""
	}
	if !token.Valid {
		return 401, nil, ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 401, nil, ""
	}
	uid := claims["sub"].(string)
	return 0, nil, uid
}
