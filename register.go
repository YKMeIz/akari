package akari

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// tokenRegister checks connected client(device)'s identification.
//
// It waits 30 second for client declaring identification.
// tokenRegister will reject connection if:
// - client keeps silence
// - client sends repeated name/token
// - client sends wrong name/token
//
// Identification declaration format:
// { "Name": "DEVICE NAME", "Token": "DEVICE TOKEN"}
func (c Core) tokenRegister(h *hub, w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	token := make(chan User)
	defer close(token)
	go func() {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		token <- readToken(string(message))
	}()

	select {
	case res := <-token:
		if res.ifRepeat(h) {
			u := User{Name: res.Name, Token: res.Token}
			if u.IsUser() {
				go response(conn, REGISTEROK)
				client := &client{hub: h, conn: conn, send: make(chan []byte, 256), user: u}
				client.hub.register <- client
				go client.writePump()
				go client.readPump(h, c)
			} else {
				go response(conn, formatErrInfo(REGISTERER))
			}
		} else {
			go response(conn, formatErrInfo(REGISTEROL))
		}
	case <-time.After(time.Second * 30):
		go response(conn, formatErrInfo(REGISTERTO))
	}
}

// readToken reads message and return in type "User".
func readToken(message string) User {
	var m User
	dec := json.NewDecoder(strings.NewReader(message))
	if err := dec.Decode(&m); err != nil {
		log.Println(err)
	}
	return m
}

// ifRepeat checks if name/token is repeated.
// It compares with online devices.
func (u User) ifRepeat(h *hub) bool {
	for client := range h.clients {
		if u == client.user {
			return false
		}
	}
	return true
}


// response makes a response to device(client) trying to make a connection.
func response(conn *websocket.Conn, reason string) {
	r := strings.NewReader(reason)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Println(err)
	}
	if reason != REGISTEROK {
		conn.Close()
	}
}
