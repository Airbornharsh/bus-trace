package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/http-server/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin
func BusCreate(c *gin.Context) {
	tx := db.DB.Begin()
	code, tempUser, err := helpers.TokenToUid(c)
	if code != 0 && err != nil {
		tx.Rollback()
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	var checkBus *models.Bus
	fmt.Println(checkBus)
	result := tx.Model(&models.Bus{}).Where("id = ?", tempUser.BusID).Find(&checkBus)
	if result.Error != gorm.ErrRecordNotFound {
		fmt.Println(checkBus)
		tx.Rollback()
		c.JSON(409, gin.H{
			"message": "Bus Already Registred for this user",
		})
		return
	}

	var bus *models.Bus
	err = c.ShouldBindJSON(&bus)
	if err != nil {
		tx.Rollback()
		fmt.Println("Unable to marse the Json")
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	bus.ID = uuid.New().String()
	result = tx.Model(&models.Bus{}).Create(&bus)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("Error in Creating", result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	result = tx.Model(&models.User{}).Where("id = ?", tempUser.ID).Update("bus_id", bus.ID).Update("bus_owner", true)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("Error in Creating", result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": "Bus Created",
		"bus":     bus,
	})
}

func SearchBus(c *gin.Context) {
	tx := db.DB.Begin()
	code, user, _ := helpers.TokenToUid(c)
	if user == nil {
		tx.Rollback()
		c.JSON(code, gin.H{
			"message": "Error While Parsing the Data",
		})
		return
	}
	search := c.Query("search")
	lat, latOk := c.Params.Get("lat")
	long, longOk := c.Params.Get("long")
	if !latOk || !longOk {
		tx.Rollback()
		c.JSON(500, gin.H{
			"message": "No Lat or No Long",
		})
		return
	}

	var buses []*models.Bus
	if err := tx.Where("ST_DWithin(ST_MakePoint(?, ?)::geography, ST_MakePoint(long, lat)::geography, ?)", long, lat, 5000).Where("LOWER(name) LIKE ?", "%"+search+"%").Find(&buses).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error in Getting List", err.Error())
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": "Bus List",
		"buses":   buses,
	})
}
