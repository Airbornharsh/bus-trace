package db

import (
	"fmt"
	"os"

	"github.com/airbornharsh/bus-trace/http-server/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbInit() {
	host := os.Getenv("HOST")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("PASSWORD")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	} else {
		DB = db
		fmt.Println("Connected")
		db.AutoMigrate(&models.User{})
		db.AutoMigrate(&models.Bus{})
	}
}
