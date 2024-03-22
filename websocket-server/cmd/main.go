package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var busClients = make(map[string][]*websocket.Conn)
var clients = make(map[*websocket.Conn]bool)
var busOwner = make(map[string]*websocket.Conn)

type BusData struct {
	BusId string `json:"busId"`
	Lat   string `json:"lat"`
	Long  string `json:"long"`
	Index int    `json:"index"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	r.GET("/ws/bus/:busId", func(c *gin.Context) {
		busId, ok := c.Params.Get("busId")
		if !ok {
			fmt.Println("Didn't got the Bus Id")
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			fmt.Println(err.Error())
		}

		clients[conn] = true
		_, has := busOwner[busId]
		if has {
			fmt.Println("Already There")
			conn.WriteMessage(1, []byte("Bus Has Already Registered"))
			conn.Close()
			delete(clients, conn)
			return
		} else {
			defer delete(busOwner, busId)
			defer removeBusConns(busId)
			defer delete(busClients, busId)
		}
		busOwner[busId] = conn
		busClients[busId] = append(busClients[busId], conn)
		conn.WriteMessage(1, []byte("Bus Added"))

		for {
			var data BusData
			if !clients[conn] {
				break
			}
			err = conn.ReadJSON(&data)
			if err != nil {
				fmt.Println("Error in Reading Json", err.Error())
			}
			data.BusId = busId
			for i, client := range busClients[busId] {
				data.Index = i + 1
				if !clients[conn] {
					continue
				}
				client.WriteJSON(data)
			}
		}
	})

	r.GET("/ws/user/:busId", func(c *gin.Context) {
		busId, ok := c.Params.Get("busId")
		if !ok {
			fmt.Println("Didn't get the Bus Id")
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			fmt.Println(err.Error())
		}

		busConn, has := busOwner[busId]
		if !has {
			defer conn.Close()
			conn.WriteMessage(1, []byte("Bus is not Yet Connected"))
			return
		} else {
			busConn.WriteMessage(1, []byte("New User is Tracking"))
		}

		clients[conn] = true
		busClients[busId] = append(busClients[busId], conn)
		conn.WriteMessage(1, []byte("Connected to the Bus"))

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(clients, conn)
				var toRemove int
				for i, v := range busClients[busId] {
					if v == conn {
						toRemove = i
					}
				}
				busClients[busId] = append(busClients[busId][:toRemove], busClients[busId][toRemove+1:]...)
				busConn.WriteMessage(1, []byte("User Discounted"))
				break
			}
		}
	})

	r.Run(":8000")
}

func removeBusConns(busId string) {
	for _, client := range busClients[busId] {
		client.WriteMessage(1, []byte("Bus Discounted"))
		client.Close()
		delete(clients, client)
	}
}
