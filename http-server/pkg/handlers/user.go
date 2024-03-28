package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/http-server/pkg/models"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	tx := db.DB.Begin()
	code, uid, err := helpers.TokenToUid2(c)
	if code != 0 && err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	var user *models.User
	err = c.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println("Unable to marse the Json")
	}

	user.ID = uid
	result := tx.Model(&models.User{}).Create(&user)

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
		"message": "User Created",
	})
}

func GetUser(c *gin.Context) {
	tx := db.DB.Begin()
	code, tempUser, err := helpers.TokenToUid(c)
	if code != 0 && err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	var user *models.User
	result := tx.Model(&models.User{}).First(&user, "id = ?", tempUser.ID)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("Error in Getting", result.Error)
		c.JSON(500, gin.H{
			"message": result.Error,
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": "Got the User",
		"user":    user,
	})
}
