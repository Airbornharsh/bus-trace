package main

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error in Loading Env File")
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	db.DBInit()
	routes.RouteInit(r)

	r.Run(":8000")
}
