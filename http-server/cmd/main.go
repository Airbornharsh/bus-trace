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
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization", "content-type"}
	r.Use(cors.New(config))

	db.DbInit()

	routes.RouteInit(r)

	go func() {
		r.Run(":8001")
	}()
	select {}
}
