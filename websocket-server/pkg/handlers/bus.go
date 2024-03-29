package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func BusSocket(c *gin.Context) {
	token, tokenOk := c.Params.Get("token")
	if !tokenOk {
		fmt.Println("Didn't got the Token")
		return
	}
	user, err := helpers.TokenToUid(token)
	if err != nil && user.ID == "" {
		fmt.Println("Parsing the Token Failed")
		return
	}
	userId := user.ID
	busId := *user.BusID
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		fmt.Println(err.Error())
	}
	websocket.UserClient[userId] = conn
	websocket.Clients[userId] = true
	_, has := websocket.BusOwner[busId]
	stopUploadCh := make(chan struct{})
	stopReadCh := make(chan struct{})
	if has {
		fmt.Println("Already There")
		conn.WriteMessage(1, []byte("Status: Bus Already Added"))
		// conn.Close()
		// delete(websocket.Clients, userId)
		// return
	} else {
		go helpers.UpdateBusLocationDB(busId, stopUploadCh)
		defer delete(websocket.BusOwner, busId)
		defer helpers.RemoveBusConns(busId)
		defer delete(websocket.BusClients, busId)
	}
	websocket.BusOwner[busId] = userId
	if clients, ok := websocket.BusClients[busId]; !ok {
		websocket.BusClients[busId] = []string{userId}
	} else {
		var found bool
		for _, id := range clients {
			if id == userId {
				found = true
				break
			}
		}
		if !found {
			websocket.BusClients[busId] = append(websocket.BusClients[busId], userId)
		}
	}
	websocket.BusLocation[busId] = types.Location{
		Lat:  0,
		Long: 0,
	}
	conn.WriteMessage(1, []byte("Status: Bus Added"))

	go func(stopCh <-chan struct{}) {
		defer fmt.Println("Goroutine exited.")
		for {
			select {
			case <-stopCh:
				fmt.Println("Stop signal received. Exiting goroutine.")
				return
			default:
				var data types.BusResponse
				if !websocket.Clients[userId] {
					break
				}
				err = conn.ReadJSON(&data)
				if err != nil {
					fmt.Println("Error in Reading Json", err.Error())
					break
				}
				if data.Which == "busData" {
					helpers.BusDataRes(userId, busId, data)
				} else if data.Which == "busClose" {
					helpers.BusCloseRes(busId, conn, data, stopReadCh, stopUploadCh)
				}
			}
		}
	}(stopReadCh)

	select {}
}
