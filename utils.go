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
	"strings"
)

type check struct {
	source      string
	destination []string
	hub         *hub
}

func checkDevice(h *hub, s string, d []string) error {
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
	u := User{Token: c.source}
	if !u.IsUser() {
		return errors.New("Source does not appear.")
	}
	for i := 0; i < len(c.destination); i++ {
		u := User{Token: c.destination[i]}
		if !u.IsUser() {
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
	if !c.isOnline(c.source) {
		u := &User{Token: c.source}
		u.UserCompletion()
		return errors.New("Source " + u.Name + " is offline.")
	}
	for i := 0; i < len(c.destination); i++ {
		if !c.isOnline(c.destination[i]) {
			u := &User{Token: c.destination[i]}
			u.UserCompletion()
			destinationName = destinationName + u.Name + " "
		}
	}
	if destinationName != "" {
		return errors.New("Destination " + destinationName + "does not appear.")
	}
	return nil
}

// isOnline checks if given destination is online.
func (c check) isOnline(token string) bool {
	for client := range c.hub.clients {
		if client.user.Token == token {
			return true
		}
	}
	return false
}

// ReadMessage reads a string of Message content, and returns in Message type.
func ReadMessage(msg string) (Message, error) {
	var m Message
	dec := json.NewDecoder(strings.NewReader(msg))
	if err := dec.Decode(&m); err == nil && messageCheck(m) {
		return m, nil
	}
	return m, errors.New(MESSAGEERR)
}

func messageCheck(m Message) bool {
	if m.Source != "" && len(m.Destination) != 0 {
		return true
	}
	return false
}
