package routes

import (
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func RouteInit(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	ws := r.Group("/ws")

	ws.GET("/bus/:userId/:busId", handlers.BusSocket)
	ws.GET("/user/:userId/:busId", handlers.UserSocket)
}
