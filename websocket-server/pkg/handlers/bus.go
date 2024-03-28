package handlers

import (
	"fmt"
	"time"

	"github.com/airbornharsh/bus-trace/websocket-server/internal/db"
	"github.com/airbornharsh/bus-trace/websocket-server/internal/websocket"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/helpers"
	"github.com/airbornharsh/bus-trace/websocket-server/pkg/models"
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
	if has {
		fmt.Println("Already There")
		conn.WriteMessage(1, []byte("Status: Bus Already Added"))
		// conn.Close()
		// delete(websocket.Clients, userId)
		// return
	} else {
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

	go func() {
		for {
			time.Sleep(60 * time.Second)
			loc, ok := websocket.BusLocation[busId]
			if ok {
				if result := db.DB.Model(&models.Bus{}).Where("id = ?", "d6e10f57-50ea-434a-a2ed-ad6a7a7205c1").Updates(map[string]interface{}{
					"lat":  loc.Lat,
					"long": loc.Long,
				}); result.Error != nil {
					fmt.Println("Update Bus Failed")
				}
			}
		}
	}()

	go func() {
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
			websocket.BusLocation[busId] = types.Location{
				Lat:  data.Lat,
				Long: data.Long,
			}
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
	}()

	select {}
}
