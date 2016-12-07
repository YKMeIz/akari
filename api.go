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
func (c Core) sendToWebsocket(g *gin.Context, hub *Hub) error {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		log.Println(err)
		return errors.New("internal error.")
	}
	msg, err := readMessage(string(body))
	if err != nil {
		return err
	}
	b, err := json.Marshal(msg.Data)
	if err != nil {
		log.Println(err)
		return errors.New("internal error.")
	}

	switch msg.Destination[0] {
	case "BROADCAST":
		hub.broadcast <- b
		return nil
	case "HANDLERFUNC":
		return runHandlerFunc(c.Event[msg.Destination[1]]) //handle err
	case "PUSHBULLET":
		if len(msg.Destination) != 1 {
			if err := sendToWebsocketDefault(hub, msg, b); err != nil {
				return err
			}
		} else {
			return makePushbulletPush(msg.Data)
		}
	default:
		return sendToWebsocketDefault(hub, msg, b)
	}
	return nil
}

func makePushbulletPush(data map[string]string) error {
	p := PushbulletPush{
		pushType:    data["Type"],
		title:       data["Title"],
		body:        data["Body"],
		accessToken: data["AccessToken"]}
	if err := p.Push(); err != nil {
		return errors.New("Fail to send Pushbullet notification.")
	}
	return nil
}

func sendToWebsocketDefault(hub *Hub, msg Message, b []byte) error {
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
