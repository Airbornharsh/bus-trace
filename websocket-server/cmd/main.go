package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var busClients = make(map[string][]string)
var clients = make(map[string]bool)
var busOwner = make(map[string]string)
var clientLocation = make(map[string]Location)
var userClient = make(map[string]*websocket.Conn)

// Bus Connected = 1
// Bus Disconnected = 2
// User Connected = 3
// User Disconnected = 4
// Bus Already Added = 5

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type BusData struct {
	BusId string  `json:"busId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	Index int     `json:"index"`
}

type UserData struct {
	UserId string  `json:"userId"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
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

	r.GET("/ws/bus/:userId/:busId", func(c *gin.Context) {
		busId, busOk := c.Params.Get("busId")
		if !busOk {
			fmt.Println("Didn't got the Bus Id")
		}
		userId, userOk := c.Params.Get("userId")
		if !userOk {
			fmt.Println("Didn't got the User Id")
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			fmt.Println(err.Error())
		}
		userClient[userId] = conn
		clients[userId] = true
		_, has := busOwner[busId]
		if has {
			fmt.Println("Already There")
			conn.WriteMessage(1, []byte("Status: Bus Already Added"))
			conn.Close()
			delete(clients, userId)
			return
		} else {
			defer delete(busOwner, busId)
			defer removeBusConns(busId)
			defer delete(busClients, busId)
		}
		busOwner[busId] = userId
		busClients[busId] = append(busClients[busId], userId)
		conn.WriteMessage(1, []byte("Status: Bus Added"))

		for {
			var data BusData
			if !clients[userId] {
				break
			}
			err = conn.ReadJSON(&data)
			if err != nil {
				fmt.Println("Error in Reading Json", err.Error())
				break
			}
			data.BusId = busId
			for i, client := range busClients[busId] {
				data.Index = i + 1
				if !clients[client] {
					continue
				}
				clientLoc := clientLocation[client]
				clientLocation[userId] = Location{
					Lat:  data.Lat,
					Long: data.Long,
				}
				if isInside(data.Lat, data.Long, clientLoc.Lat, clientLoc.Long) {
					// if client != userId {
					// 	conn.WriteJSON(UserData{
					// 		UserId: userId,
					// 		Lat:    clientLoc.Lat,
					// 		Long:   clientLoc.Long,
					// 	})
					// }
					con, has := userClient[client]
					if has {
						con.WriteJSON(data)
					}
				}
			}
		}
	})

	r.GET("/ws/user/:userId/:busId", func(c *gin.Context) {
		busId, busOk := c.Params.Get("busId")
		if !busOk {
			fmt.Println("Didn't get the Bus Id")
		}
		userId, userOk := c.Params.Get("userId")
		if !userOk {
			fmt.Println("Didn't got the User Id")
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			fmt.Println(err.Error())
		}

		busUserId, has := busOwner[busId]
		if !has {
			defer conn.Close()
			conn.WriteMessage(1, []byte("Status: Bus is not Connected"))
			return
		} else {
			busConn, has := userClient[busUserId]
			if !has {
				fmt.Println("Error in getting the bus Conn")
			}
			busConn.WriteMessage(1, []byte("Status: User is Tracking"))
		}

		userClient[userId] = conn
		clients[userId] = true
		busClients[busId] = append(busClients[busId], userId)
		busUserId, busUserIdHas := busOwner[busId]
		if !busUserIdHas {
			fmt.Println("Not Found Bus User Id")
		}
		busLocation, busLocationHas := clientLocation[busUserId]
		if !busLocationHas {
			fmt.Println("No Bus Location Found")
		}
		conn.WriteJSON(BusData{
			BusId: busId,
			Lat:   busLocation.Lat,
			Long:  busLocation.Long,
			Index: 1,
		})
		conn.WriteMessage(1, []byte("Status: Connected to Bus"))

		for {
			var loc Location
			err := conn.ReadJSON(&loc)
			fmt.Println(loc)
			busConn, has := userClient[busUserId]
			if !has {
				fmt.Println("Error in getting the bus Conn")
			}
			if err != nil {
				delete(clients, userId)
				delete(clientLocation, userId)
				var toRemove int
				for i, v := range busClients[busId] {
					if v == userId {
						toRemove = i
					}
				}
				busClients[busId] = append(busClients[busId][:toRemove], busClients[busId][toRemove+1:]...)
				busConn.WriteMessage(1, []byte("Status: User Disconnected"))
				break
			}
			clientLocation[userId] = Location{
				Lat:  loc.Lat,
				Long: loc.Long,
			}
			busConn.WriteJSON(UserData{
				UserId: userId,
				Lat:    loc.Lat,
				Long:   loc.Long,
			})
		}
	})

	r.Run(":8000")
}

func removeBusConns(busId string) {
	for _, userId := range busClients[busId] {
		client, has := userClient[userId]
		if !has {
			fmt.Println("Error in getting the bus Conn")
		}
		client.WriteMessage(1, []byte("Status: Bus Disconnected"))
		client.Close()
		_, locationHas := clientLocation[userId]
		if locationHas {
			delete(clientLocation, userId)
		}
		delete(clients, userId)
	}
}

func isInside(lat1, lon1, lat2, lon2 float64) bool {
	const radiusOfEarth = 6371 // Radius of Earth in kilometers

	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLon := lon2Rad - lon1Rad
	dLat := lat2Rad - lat1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := radiusOfEarth * c

	if distance > 10000 {
		return false
	} else {
		return true
	}
}
