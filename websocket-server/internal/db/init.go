package db

import (
	"fmt"
	"os"

	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInit() {
	host := os.Getenv("HOST")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("PASSWORD")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		DB = db
		fmt.Println("Connected")
		db.AutoMigrate(&models.User{})
		db.AutoMigrate(&models.Bus{})
		// db.Migrator().DropTable(&models.User{}, &models.Bus{})
	}
}
