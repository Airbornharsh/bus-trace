package routes

import (
	"github.com/airbornharsh/bus-trace/http-server/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func RouteInit(r *gin.Engine) {
	r.GET("", func(c *gin.Context) {
		c.JSON(
			200, gin.H{
				"message": "Hello world",
			},
		)
	})

	busRoutes := r.Group("bus")
	userRoutes := r.Group("user")

	userRoutes.POST("", handlers.CreateUser)
	userRoutes.GET("", handlers.GetUser)
	userRoutes.GET("/:userId", handlers.GetUserById)
	busRoutes.POST("", handlers.BusCreate)
	busRoutes.GET("/:lat/:long", handlers.SearchBus)
}
