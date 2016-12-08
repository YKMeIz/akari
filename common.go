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

import ()

// Core is the framework's instance, it contains all configuration settings.
// Create an instance of Core, by using New().
type Core struct {
	// Server's domain name.
	Domain string

	// Port number listened by Akari Message Framework.
	Port string

	// Domain's certificate chain.
	CertChain string

	// Domain's privatekey.
	CertKey string

	// Relative path for handling HTTP POST requests.
	MessageRelativePath string

	// Relative path for providing websocket service.
	WebsocketRelativePath string

	// Path to SQLite database file.
	DatabasePath string

	// Event is a map points to your custom functions.
	Event map[string]HandlerFunc
	// Examples could be found on https://github.com/nrechn/akari .
}

type HandlerFunc func(*Message) error

// Message defines "Unified Message Format".
// Examples could be found on https://github.com/nrechn/akari .
type Message struct {
	// Message sender's token.
	Source string

	// Message receiver's token.
	Destination []string

	// Message content.
	Data map[string]string
}

// PushbulletPush defines Pushbullet's push action (notification).
type PushbulletPush struct {
	// Type of the push, one of "note", "file", "link".
	// Akari currently only supports type "note".
	PushType string

	// Title of the push, used for all types of pushes.
	Title string

	// Body of the push, used for all types of pushes.
	Body string

	// Access Token of your account.
	AccessToken string
}

// Constant strings for return request/message status.
const (
	REGISTEROK string = `{"Status": "ok!"}`
	REGISTERER string = `Wrong user name or token.`
	REGISTEROL string = `Your device is already online.`
	REGISTERTO string = `Id authentication is required.`
	PBPUSHNERR string = `Fail to send Pushbullet notification.`
	HANDLERFER string = `An error occured on running custom handler function.`
	MESSAGEERR string = `Missing source or destination.`
)

var (
	DatabasePath string
)

func formatErrInfo(errInfo string) string {
	return `{"Status": "error! ` + errInfo + `"}`
}
