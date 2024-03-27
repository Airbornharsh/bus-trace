package helpers

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenToUid(c *gin.Context) (int, error, string) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return 401, nil, ""
	}
	auth2 := strings.Split(authorization, " ")
	if len(auth2) < 2 {
		return 401, nil, ""
	}
	tokenString := auth2[1]
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
