package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func BusSocket(c *gin.Context) {
	busId, busOk := c.Params.Get("busId")
	if !busOk {
		fmt.Println("Didn't got the Bus Id")
	}
	userId, userOk := c.Params.Get("userId")
	if !userOk {
		fmt.Println("Didn't got the User Id")
	}
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		fmt.Println(err.Error())
	}
	websocket.UserClient[userId] = conn
	websocket.Clients[userId] = true
	_, has := websocket.BusOwner[busId]
	if has {
		fmt.Println("Already There")
		conn.WriteMessage(1, []byte("Status: Bus Already Added"))
		conn.Close()
		delete(websocket.Clients, userId)
		return
	} else {
		defer delete(websocket.BusOwner, busId)
		defer helpers.RemoveBusConns(busId)
		defer delete(websocket.BusClients, busId)
	}
	websocket.BusOwner[busId] = userId
	websocket.BusClients[busId] = append(websocket.BusClients[busId], userId)
	conn.WriteMessage(1, []byte("Status: Bus Added"))

	for {
		var data types.BusData
		if !websocket.Clients[userId] {
			break
		}
		err = conn.ReadJSON(&data)
		if err != nil {
			fmt.Println("Error in Reading Json", err.Error())
			break
		}
		data.BusId = busId
		for i, client := range websocket.BusClients[busId] {
			data.Index = i + 1
			if !websocket.Clients[client] {
				continue
			}
			clientLoc := websocket.ClientLocation[client]
			websocket.ClientLocation[userId] = types.Location{
				Lat:  data.Lat,
				Long: data.Long,
			}
			if helpers.IsInside(data.Lat, data.Long, clientLoc.Lat, clientLoc.Long) {
				// if client != userId {
				// 	conn.WriteJSON(UserData{
				// 		UserId: userId,
				// 		Lat:    clientLoc.Lat,
				// 		Long:   clientLoc.Long,
				// 	})
				// }
				con, has := websocket.UserClient[client]
				if has {
					con.WriteJSON(data)
				}
			}
		}
	}
}
