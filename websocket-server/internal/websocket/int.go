package websocket

import (
	"net/http"

	types "github.com/airbornharsh/bus-trace/websocket-server/pkg/Types"
	"github.com/gorilla/websocket"
)

var BusClients = make(map[string][]string)
var Clients = make(map[string]bool)
var BusOwner = make(map[string]string)
var ClientLocation = make(map[string]types.Location)
var UserClient = make(map[string]*websocket.Conn)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
