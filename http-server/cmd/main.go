package main

import (
	"fmt"
	"net/http"

	"github.com/airbornharsh/bus-trace/http-server/internal/db"
	"github.com/airbornharsh/bus-trace/http-server/pkg/routes"
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

	r.Use(CorsMiddleware())

	db.DbInit()

	r.HEAD("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/", func(c *gin.Context) {
		fmt.Println(c)
		fmt.Println("Hello World POST")
		c.JSON(200, gin.H{
			"message": "Welcome to Bus Trace Http Backend POST",
		})
	})

	routes.RouteInit(r)

	fmt.Println("Routes")

	r.Run("0.0.0.0:8001")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin with the specified methods
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Continue with the next middleware/handler
		c.Next()
	}
}
