package main

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error in Loading Env File")
	}

	db.DbInit()

	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	r.POST("/", func(c *gin.Context) {
		fmt.Println(c)
		fmt.Println("Hello World POST")
		c.JSON(200, gin.H{
			"message": "Welcome to Bus Trace Http Backend POST",
		})
	})

	routes.RouteInit(r)

	fmt.Println("Server Started at http://localhost:8001")
	r.Run(":8001")
}
