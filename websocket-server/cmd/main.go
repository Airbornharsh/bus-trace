package main

import (
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	routes.RouteInit(r)

	go func() {
		r.Run(":8000")
	}()
	select {}
}
