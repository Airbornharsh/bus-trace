package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers/mutex"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

func UserSocket(c *gin.Context) {
	busId, busOk := c.Params.Get("busId")
	if !busOk {
		fmt.Println("Didn't get the Bus Id")
	}
	token, tokenOk := c.Params.Get("token")
	if !tokenOk {
		fmt.Println("Didn't got the Token")
	}
	user, err := helpers.TokenToUid(token)
	if err != nil && user.ID == "" {
		fmt.Println("Parsing the Token Failed")
		return
	}
	userId := user.ID
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		fmt.Println(err.Error())
	}

	busUserId, has := websocket.BusOwner[busId]
	if !has {
		defer conn.Close()
		// conn.WriteMessage(1, []byte("Status: Bus is not Connected"))
		conn.WriteJSON(types.UserResponse{
			UserId: userId,
			UserMessage: types.UserMessage{
				Message: "Bus is not Connected",
			},
			Which: "userMessage",
		})
		return
	} else {
		var busConn *ws.Conn
		var has bool
		for {
			if mutex.UserClientMutex {
				mutex.UserClientMutex = false
				busConn, has = websocket.UserClient[busUserId]
				mutex.UserClientMutex = true
				break
			}
		}
		if !has {
			fmt.Println("Error in getting the bus Conn")
		}
		// busConn.WriteMessage(1, []byte("Status: User is Tracking"))
		busConn.WriteJSON(types.BusResponse{
			BusId: busId,
			BusMessage: types.BusMessage{
				Message: "User is Tracking",
			},
			Which: "busMessage",
		})
		for {
			if mutex.UserClientMutex {
				mutex.UserClientMutex = false
				websocket.UserClient[userId] = conn
				mutex.UserClientMutex = true
				break
			}
		}
		for {
			if mutex.ClientsMutex {
				mutex.ClientsMutex = false
				websocket.Clients[userId] = true
				mutex.ClientsMutex = true
				break
			}
		}
		for {
			if mutex.BusClientsMutex {
				mutex.BusClientsMutex = false
				websocket.BusClients[busId] = append(websocket.BusClients[busId], userId)
				mutex.BusClientsMutex = true
				break
			}
		}
		busUserId, busUserIdHas := websocket.BusOwner[busId]
		if !busUserIdHas {
			fmt.Println("Not Found Bus User Id")
		}
		busLocation, busLocationHas := websocket.ClientLocation[busUserId]
		if !busLocationHas {
			fmt.Println("No Bus Location Found")
		}
		conn.WriteJSON(types.UserResponse{
			UserBusData: types.UserBusData{
				BusId: busId,
				Lat:   busLocation.Lat,
				Long:  busLocation.Long,
			},
			UserMessage: types.UserMessage{
				Message: "Connected to Bus",
			},
			Which: "userBusData&userMessage",
		})

		for {
			var data types.UserRequest
			if !websocket.Clients[userId] {
				break
			}
			err := conn.ReadJSON(&data)
			var busConn *ws.Conn
			var has bool
			for {
				if mutex.UserClientMutex {
					mutex.UserClientMutex = false
					busConn, has = websocket.UserClient[busUserId]
					mutex.UserClientMutex = true
					break
				}
			}
			if !has {
				fmt.Println("Error in getting the bus Conn")
			}
			if err != nil {
				for {
					if mutex.ClientsMutex {
						mutex.ClientsMutex = false
						delete(websocket.Clients, userId)
						mutex.ClientsMutex = true
						break
					}
				}
				for {
					if mutex.ClientLocationMutex {
						mutex.ClientLocationMutex = false
						delete(websocket.ClientLocation, userId)
						mutex.ClientLocationMutex = true
						break
					}
				}
				var tempSlice []string
				for {
					if mutex.BusClientsMutex {
						mutex.BusClientsMutex = false
						var toRemove int
						for i, v := range websocket.BusClients[busId] {
							if v == userId {
								toRemove = i
							}
						}
						tempSlice = append(websocket.BusClients[busId][:toRemove], websocket.BusClients[busId][toRemove+1:]...)
						websocket.BusClients[busId] = tempSlice
						mutex.BusClientsMutex = true
						break
					}
				}
				busConn.WriteJSON(types.BusResponse{
					BusUserList: tempSlice,
					BusMessage: types.BusMessage{
						Message: "User Disconnected",
					},
					Which: "busMessage&busUserList",
				})
				break
			}
			if data.Which == "userLocation" {
				for {
					if mutex.ClientLocationMutex {
						mutex.ClientLocationMutex = false
						websocket.ClientLocation[userId] = types.Location{
							Lat:  data.UserLocation.Lat,
							Long: data.UserLocation.Long,
						}
						mutex.ClientLocationMutex = true
						break
					}
				}
				helpers.UpdateUserLocationDB(userId, data.UserLocation.Lat, data.UserLocation.Long)
				busConn.WriteJSON(types.BusResponse{
					BusId:       busId,
					BusUserList: websocket.BusClients[busId],
					BusUserLocation: types.BusUserLocation{
						UserId: userId,
						Lat:    data.UserLocation.Lat,
						Long:   data.UserLocation.Long,
					},
					Which: "busUserLocation&busUserList",
				})
			}
		}
	}
}
