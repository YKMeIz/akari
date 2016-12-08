package akari

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump(h *hub, co Core) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg, err := ReadMessage(string(message))
		if err != nil {
			c.send <- []byte(formatErrInfo(err.Error()))
		} else {
			b, err := json.Marshal(msg.Data)
			if err != nil {
				log.Println(err)
				break
			}
			switch msg.Destination[0] {
			case "BROADCAST":
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				c.hub.broadcast <- message
			case "HANDLERFUNC":
				if err := msg.runHandlerFunc(co.Event[msg.Destination[1]]); err != nil {
					c.send <- []byte(`{"Status":"error! ` + err.Error() + `"}`)
				} //handle err
			case "PUSHBULLET":
				if len(msg.Destination) != 1 {
					c.sendToWebsocketReadPump(h, msg, b)
				} else {
					if err := makePushbulletPush(msg.Data); err != nil {
						c.send <- []byte(formatErrInfo(PBPUSHNERR))
					} //handle err
				}
			default:
				c.sendToWebsocketReadPump(h, msg, b)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) sendToWebsocketReadPump(h *hub, msg Message, b []byte) {
	r := checkDevice(h, msg.Source, msg.Destination)
	if r != nil {
		c.send <- []byte(`{"Status":"error! ` + r.Error() + `"}`)
	} else {
		for i := 0; i < len(msg.Destination); i++ {
			for client := range h.clients {
				if client.user.Token == msg.Destination[i] {
					client.send <- b
				}
			}
		}
	}
}
