package handlers

import (
	"fmt"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers/mutex"
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
	for {
		if mutex.BusLocationMutex {
			mutex.BusLocationMutex = false
			websocket.UserClient[userId] = conn
			mutex.BusLocationMutex = true
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
	_, has := websocket.BusOwner[busId]
	stopUploadCh := make(chan struct{})
	stopReadCh := make(chan struct{})
	if has {
		fmt.Println("Already There")
		// conn.WriteMessage(1, []byte("Status: Bus Already Added"))
		// conn.Close()
		// delete(websocket.Clients, userId)
		// return
		conn.WriteJSON(types.BusResponse{
			BusId: busId,
			BusMessage: types.BusMessage{
				Message: "Bus Already Added",
			},
			Which: "busMessage",
		})
	} else {
		go helpers.UpdateBusLocationDB(busId, stopUploadCh)
		defer func() {
			for {
				if mutex.BusOwnerMutex {
					mutex.BusOwnerMutex = false
					delete(websocket.BusOwner, busId)
					mutex.BusOwnerMutex = true
					break
				}
			}
		}()
		defer helpers.RemoveBusConns(busId)
		defer func() {
			for {
				if mutex.BusClientsMutex {
					mutex.BusClientsMutex = false
					delete(websocket.BusClients, busId)
					mutex.BusClientsMutex = true
					break
				}
			}
		}()
	}
	for {
		if mutex.BusOwnerMutex {
			mutex.BusOwnerMutex = false
			websocket.BusOwner[busId] = userId
			mutex.BusOwnerMutex = true
			break
		}
	}
	if clients, ok := websocket.BusClients[busId]; !ok {
		for {
			if mutex.BusClientsMutex {
				mutex.BusClientsMutex = false
				websocket.BusClients[busId] = []string{userId}
				mutex.BusClientsMutex = true
				break
			}
		}
	} else {
		var found bool
		for _, id := range clients {
			if id == userId {
				found = true
				break
			}
		}
		if !found {
			for {
				if mutex.BusClientsMutex {
					mutex.BusClientsMutex = false
					websocket.BusClients[busId] = append(websocket.BusClients[busId], userId)
					mutex.BusClientsMutex = true
					break
				}
			}
		}
	}
	for {
		if mutex.BusLocationMutex {
			mutex.BusLocationMutex = false
			websocket.BusLocation[busId] = types.Location{
				Lat:  0,
				Long: 0,
			}
			mutex.BusLocationMutex = true
			break
		}
	}
	// conn.WriteMessage(1, []byte("Status: Bus Added"))
	conn.WriteJSON(types.BusResponse{
		BusId: busId,
		BusMessage: types.BusMessage{
			Message: "Bus Added",
		},
		Which: "busMessage",
	})

	go func(stopCh <-chan struct{}) {
		defer fmt.Println("Goroutine exited.")
	DataRead:
		for {
			select {
			case <-stopCh:
				fmt.Println("Stop signal received. Exiting goroutine.")
				return
			default:
				var data types.BusRequest
				if !websocket.Clients[userId] {
					break
				}
				err = conn.ReadJSON(&data)
				if err != nil {
					fmt.Println("Error in Reading Json", err.Error())
					close(stopReadCh)
					close(stopUploadCh)
					break DataRead
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
