package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/gin-gonic/gin"
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
		conn.WriteMessage(1, []byte("Status: Bus is not Connected"))
		return
	} else {
		busConn, has := websocket.UserClient[busUserId]
		if !has {
			fmt.Println("Error in getting the bus Conn")
		}
		busConn.WriteMessage(1, []byte("Status: User is Tracking"))

		websocket.UserClient[userId] = conn
		websocket.Clients[userId] = true
		websocket.BusClients[busId] = append(websocket.BusClients[busId], userId)
		busUserId, busUserIdHas := websocket.BusOwner[busId]
		if !busUserIdHas {
			fmt.Println("Not Found Bus User Id")
		}
		busLocation, busLocationHas := websocket.ClientLocation[busUserId]
		if !busLocationHas {
			fmt.Println("No Bus Location Found")
		}
		conn.WriteJSON(types.BusData{
			BusId: busId,
			Lat:   busLocation.Lat,
			Long:  busLocation.Long,
			Index: 1,
		})
		conn.WriteMessage(1, []byte("Status: Connected to Bus"))
		busConn.WriteJSON(websocket.BusClients[busId])

		for {
			var loc types.Location
			err := conn.ReadJSON(&loc)
			fmt.Println(loc)
			busConn, has := websocket.UserClient[busUserId]
			if !has {
				fmt.Println("Error in getting the bus Conn")
			}
			if err != nil {
				delete(websocket.Clients, userId)
				delete(websocket.ClientLocation, userId)
				var toRemove int
				for i, v := range websocket.BusClients[busId] {
					if v == userId {
						toRemove = i
					}
				}
				tempSlice := append(websocket.BusClients[busId][:toRemove], websocket.BusClients[busId][toRemove+1:]...)
				websocket.BusClients[busId] = tempSlice
				busConn.WriteMessage(1, []byte("Status: User Disconnected"))
				busConn.WriteJSON(tempSlice)
				break
			}
			websocket.ClientLocation[userId] = types.Location{
				Lat:  loc.Lat,
				Long: loc.Long,
			}
			busConn.WriteJSON(types.UserData{
				UserId: userId,
				Lat:    loc.Lat,
				Long:   loc.Long,
			})
		}
	}
}