package helpers

import (
	"fmt"
	"math"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
)

func RemoveBusConns(busId string) {
	for _, userId := range websocket.BusClients[busId] {
		client, has := websocket.UserClient[userId]
		if !has {
			fmt.Println("Error in getting the bus Conn")
		}
		client.WriteMessage(1, []byte("Status: Bus Disconnected"))
		client.Close()
		_, locationHas := websocket.ClientLocation[userId]
		if locationHas {
			delete(websocket.ClientLocation, userId)
		}
		delete(websocket.Clients, userId)
	}
}

func IsInside(lat1, lon1, lat2, lon2 float64) bool {
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
