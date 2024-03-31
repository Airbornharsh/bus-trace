package helpers

import (
	"fmt"
	"math"
	"time"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers/mutex"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
	ws "github.com/gorilla/websocket"
)

func BusDataRes(userId string, busId string, data types.BusRequest) {
	for {
		if mutex.BusLocationMutex {
			mutex.BusLocationMutex = false
			websocket.BusLocation[busId] = types.Location{
				Lat:  data.BusData.Lat,
				Long: data.BusData.Long,
			}
			mutex.BusLocationMutex = true
			break
		}
	}
	for i, client := range websocket.BusClients[busId] {
		data.BusData.Index = i + 1
		if !websocket.Clients[client] {
			continue
		}
		clientLoc := websocket.ClientLocation[client]
		for {
			if mutex.ClientLocationMutex {
				mutex.ClientLocationMutex = false
				websocket.ClientLocation[userId] = types.Location{
					Lat:  data.BusData.Lat,
					Long: data.BusData.Long,
				}
				mutex.ClientLocationMutex = true
				break
			}
		}
		if IsInside(data.BusData.Lat, data.BusData.Long, clientLoc.Lat, clientLoc.Long) {
			// if client != userId {
			// 	conn.WriteJSON(types.BusResponse{
			// 		BusId: busId,
			// 		BusUserLocation: types.BusUserLocation{
			// 			UserId: userId,
			// 			Lat:    clientLoc.Lat,
			// 			Long:   clientLoc.Long,
			// 		},
			// 		Which: "busUserLocation",
			// 	})
			// }
			var con *ws.Conn
			var has bool
			for {
				if mutex.UserClientMutex {
					mutex.UserClientMutex = false
					con, has = websocket.UserClient[client]
					mutex.UserClientMutex = true
					break
				}
			}
			if has {
				con.WriteJSON(types.UserResponse{
					UserId: userId,
					UserBusData: types.UserBusData{
						BusId: busId,
						Lat:   data.BusData.Lat,
						Long:  data.BusData.Long,
					},
					Which: "userBusData",
				})
			}
		}
	}
}

func BusCloseRes(busId string, conn *ws.Conn, data types.BusRequest, stopReadCh chan struct{}, stopUploadCh chan struct{}) {
	if data.BusClose.Close {
		conn.Close()
		select {
		case <-stopReadCh:
		default:
			close(stopReadCh)
		}

		// Signal to stop uploading
		select {
		case <-stopUploadCh:
		default:
			close(stopUploadCh)
		}
		for {
			if mutex.BusOwnerMutex {
				mutex.BusOwnerMutex = false
				delete(websocket.BusOwner, busId)
				mutex.BusOwnerMutex = true
				break
			}
		}
		RemoveBusConns(busId)
		for {
			if mutex.BusClientsMutex {
				mutex.BusClientsMutex = false
				delete(websocket.BusClients, busId)
				mutex.BusClientsMutex = true
				break
			}
		}
		return
	}
}

func RemoveBusConns(busId string) {
	for _, userId := range websocket.BusClients[busId] {
		var client *ws.Conn
		var has bool
		for {
			if mutex.UserClientMutex {
				mutex.UserClientMutex = false
				client, has = websocket.UserClient[userId]
				mutex.UserClientMutex = true
				break
			}
		}
		if !has {
			fmt.Println("Error in getting the bus Conn")
		}
		for {
			if mutex.ClientsMutex {
				mutex.ClientsMutex = false
				delete(websocket.Clients, userId)
				mutex.ClientsMutex = true
				break
			}
		}
		_, locationHas := websocket.ClientLocation[userId]
		if locationHas {
			for {
				if mutex.ClientLocationMutex {
					mutex.ClientLocationMutex = false
					delete(websocket.ClientLocation, userId)
					mutex.ClientLocationMutex = true
					break
				}
			}
		}
		// client.WriteMessage(1, []byte("Status: Bus Disconnected"))
		client.WriteJSON(types.UserResponse{
			UserMessage: types.UserMessage{
				Message: "Bus Disconnected",
			},
			Which: "userMessage",
		})
		client.Close()
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

func UpdateBusLocationDB(busId string, stopCh <-chan struct{}) {
	defer fmt.Println("Goroutine exited.")
	for {
		select {
		case <-stopCh:
			fmt.Println("Stop signal received. Exiting goroutine.")
			return
		case <-time.After(60 * time.Second):
			var loc types.Location
			var ok bool
			for {
				if mutex.BusLocationMutex {
					mutex.BusLocationMutex = false
					loc, ok = websocket.BusLocation[busId]
					mutex.BusLocationMutex = true
					break
				}
			}
			if ok {
				if result := db.DB.Model(&models.Bus{}).Where("id = ?", busId).Updates(map[string]interface{}{
					"lat":  loc.Lat,
					"long": loc.Long,
				}); result.Error != nil {
					fmt.Println("Update Bus Failed")
				}
			}
		}
	}
}
