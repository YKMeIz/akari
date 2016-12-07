package akari

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

// PushNoti handles notification pushing via websocket.
//
// If target device is offline, PushNoti will send notificatrion via Pushbullet
// as long as Pushbullet's token is set.
func sendToWebsocket(g *gin.Context, hub *Hub) error {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		log.Println(err)
		return errors.New("internal error.")
	}
	msg := readMessage(string(body))
	b, err := json.Marshal(msg.Data)
	if err != nil {
		log.Println(err)
		return errors.New("internal error.")
	}

	if msg.Destination[0] == "BROADCAST" {
		hub.broadcast <- b
		return nil
	} else {
		r := checkDevice(hub, msg.Source, msg.Destination)
		if r != nil {
			return r
		}
		for i := 0; i < len(msg.Destination); i++ {
			for client := range hub.clients {
				if client.token == msg.Destination[i] {
					client.send <- b
				}
			}
		}
		return nil
	}
}
