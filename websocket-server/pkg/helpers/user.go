package helpers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
)

func UpdateUserLocationDB(userId string, lat float64, long float64) {
	if result := db.DB.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"lat":  lat,
		"long": long,
	}); result.Error != nil {
		fmt.Println("Update User Failed")
	}
}
