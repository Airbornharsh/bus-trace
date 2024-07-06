package helpers

import (
	"os"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func TokenToUid(tokenString string) (*models.User, error) {
	var user *models.User
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return user, err
	}
	if !token.Valid {
		return user, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, nil
	}
	uid := claims["sub"].(string)

	result := db.DB.Model(&models.User{}).First(&user, "id = ?", uid)
	if result.Error == gorm.ErrRecordNotFound {
		return user, result.Error
	}
	return user, nil
}
