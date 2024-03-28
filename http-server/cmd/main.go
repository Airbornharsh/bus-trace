package main

import (
	"fmt"
	"net/http"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error in Loading Env File")
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	db.DbInit()

	r.HEAD("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	routes.RouteInit(r)

	fmt.Println("Routes")

	r.Run(":8001")
}
