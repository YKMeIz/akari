package akari

import (
	"log"
	"net/http"
)

// AuthMessage is utilized to verify device's identity.
type AuthMessage struct {
	Name, Token string
}

// handleWebsocket handles websocket requests from the peer.
func handleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go tokenRegister(hub, w, r, conn)
}
