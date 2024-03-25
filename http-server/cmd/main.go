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
	r := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error in Loading Env File")
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	db.DbInit()

	routes.RouteInit(r)

	r.Run(":8001")
}