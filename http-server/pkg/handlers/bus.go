package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/http-server/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BusCreate(c *gin.Context) {
	code, tempUser, err := helpers.TokenToUid(c)
	if code != 0 && err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	var bus *models.Bus
	err = c.ShouldBindJSON(&bus)
	if err != nil {
		fmt.Println("Unable to marse the Json")
	}

	bus.ID = uuid.New().String()
	result := db.DB.Create(&bus)
	if result.Error != nil {
		fmt.Println("Error in Creating", result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	result = db.DB.Model(&models.User{}).Where("id = ?", tempUser.ID).Update("bus_id", bus.ID).Update("bus_owner", true)
	if result.Error != nil {
		fmt.Println("Error in Creating", result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Bus Created",
		"bus":     bus,
	})
}
