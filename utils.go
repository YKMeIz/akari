package akari

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type check struct {
	source      string
	destination []string
	hub         *Hub
}

func checkDevice(h *Hub, s string, d []string) error {
	c := check{source: s, destination: d, hub: h}
	if err := c.isTokenAppear(); err != nil {
		return err
	}
	if err := c.isTokenOnline(); err != nil {
		return err
	}
	return nil
}

func (c check) isTokenAppear() error {
	var destinationToken string
	if !isToken(c.source) {
		return errors.New("Source does not appear.")
	}
	for i := 0; i < len(c.destination); i++ {
		if !isToken(c.destination[i]) {
			destinationToken = destinationToken + c.destination[i] + " "
		}
	}
	if destinationToken != "" {
		return errors.New("Destination " + destinationToken + "does not appear.")
	}
	return nil
}

func (c check) isTokenOnline() error {
	var destinationName string
	if !c.hub.isOnline(c.source) {
		return errors.New("Source " + getName(c.source) + " is offline.")
	}
	for i := 0; i < len(c.destination); i++ {
		if !c.hub.isOnline(c.destination[i]) {
			destinationName = destinationName + getName(c.destination[i]) + " "
		}
	}
	if destinationName != "" {
		return errors.New("Destination " + destinationName + "does not appear.")
	}
	return nil
}

// isOnline checks if given destination is online.
func (h *Hub) isOnline(token string) bool {
	for client := range h.clients {
		if client.token == token {
			return true
		}
	}
	return false
}

// readMessage reads a string of Message content, and returns in Message type.
func readMessage(msg string) Message {
	var m Message
	dec := json.NewDecoder(strings.NewReader(msg))
	if err := dec.Decode(&m); err != nil {
		log.Println(err)
	}
	return m
}
