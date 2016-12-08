// Copyright © 2016 nrechn <nrechn@gmail.com>
//
// This file is part of akari.
//
// akari is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// akari is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with akari. If not, see <http://www.gnu.org/licenses/>.
//

package akari

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func (c Core) sendToWebsocket(g *gin.Context, h *hub) error {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		log.Println(err)
		return errors.New("internal error.")
	}
	msg, err := ReadMessage(string(body))
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
		h.broadcast <- b
		return nil
	case "HANDLERFUNC":
		return msg.runHandlerFunc(c.Event[msg.Destination[1]]) //handle err
	case "PUSHBULLET":
		if len(msg.Destination) != 1 {
			if err := sendToWebsocketDefault(h, msg, b); err != nil {
				return err
			}
		} else {
			return makePushbulletPush(msg.Data)
		}
	default:
		return sendToWebsocketDefault(h, msg, b)
	}
	return nil
}

func makePushbulletPush(data map[string]string) error {
	p := PushbulletPush{
		PushType:    data["Type"],
		Title:       data["Title"],
		Body:        data["Body"],
		AccessToken: data["AccessToken"]}
	if err := p.Push(); err != nil {
		return errors.New(PBPUSHNERR)
	}
	return nil
}

func sendToWebsocketDefault(h *hub, msg Message, b []byte) error {
	r := checkDevice(h, msg.Source, msg.Destination)
	if r != nil {
		return r
	}
	for i := 0; i < len(msg.Destination); i++ {
		for client := range h.clients {
			if client.user.Token == msg.Destination[i] {
				client.send <- b
			}
		}
	}
	return nil
}
