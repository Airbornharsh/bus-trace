package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func TokenToUid(c *gin.Context) (int, *models.User, error) {
	var user *models.User
	authorization := c.Request.Header.Get("Authorization")
	fmt.Println("Data")
	if authorization == "" {
		return 401, user, nil
	}
	auth2 := strings.Split(authorization, " ")
	if len(auth2) < 2 {
		return 401, user, nil
	}
	tokenString := auth2[1]
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		fmt.Println(tokenString)
		fmt.Println(err)
		return 403, user, err
	}
	if !token.Valid {
		return 401, user, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 401, user, nil
	}
	uid := claims["sub"].(string)

	result := db.DB.Model(&models.User{}).First(&user, "id = ?", uid)
	if result.Error == gorm.ErrRecordNotFound {
		return 404, user, result.Error
	}
	return 0, user, nil
}

func TokenToUid2(c *gin.Context) (int, string, error) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return 401, "", nil
	}
	auth2 := strings.Split(authorization, " ")
	if len(auth2) < 2 {
		return 401, "", nil
	}
	tokenString := auth2[1]
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 403, "", err
	}
	if !token.Valid {
		return 401, "", nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 401, "", nil
	}
	uid := claims["sub"].(string)
	return 0, uid, nil
}
