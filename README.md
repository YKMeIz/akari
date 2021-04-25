<img align="right" src="https://raw.githubusercontent.com/YKMeIz/akari/main/logo.png">

# Akari Message Framework
[![License](https://img.shields.io/badge/license-GPL--3.0-red.svg?style=flat-square)](https://github.com/YKMeIz/akari/blob/main/LICENSE)
[![Go Walker](https://img.shields.io/badge/Go%20Walker-documentation-green.svg?style=flat-square)](http://gowalker.org/github.com/YKMeIz/akari)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/YKMeIz/akari)
[![Go Report Card](https://goreportcard.com/badge/github.com/YKMeIz/akari?style=flat-square)](https://goreportcard.com/report/github.com/YKMeIz/akari)


Akari is a message framework written in Go (Golang). It follows KISS (Keep it simple, stupid) design principle, and is designed for IoT communication and notification push from *nix side to any device.

## Quickstart

#### Download and install source

```sh
$ go get github.com/YKMeIz/akari
```

#### Create a user
> Assume you have installed SQLite in your host system.

In the example below, we want to create new user `Akari`, and output new token for `Akari`:

```go
package main

import (
	"fmt"
	"github.com/YKMeIz/akari"
)

func main() {
	// Create a database file
	akari.InitDatabase("/tmp/data.db")

	// Create a new user
	u := akari.User{Name: "Akari"}
	name, token, err := u.RegisterUser()
	if err != nil {
		panic(err)
	}
	fmt.Println("Create User: " + name)
	fmt.Println(name + "'s token is: " + token)
}
```

Run the Go file you just write (e.g. createUser.go). The output should look like following:

```sh
$ go run createUser.go
Create User: Akari
Akari's token is: f6283b29169cf8c1e84bf23cf86772fb
```

#### Run Akari based server

Start Akari based server is quite simple. Only having a few basic settings, it will handle all dirty works. Here is an example:

```go
package main

import (
	"fmt"
	"github.com/YKMeIz/akari"
)

func main() {
	c := akari.New()
	c.DatabasePath = "/tmp/data.db"

	// Custom function
	c.Event["PRINT"] = eventPrint

	// Listen and serve on IPAddress:8080
	c.Run()
}

// Custom function
func eventPrint(m *akari.Message) error {
	fmt.Println(m.Data["customFunc"])
	return nil
}
```

Run the Go file you just write (e.g. akariServer.go). The output should look like following:

```sh
$ go run akariServer.go

Server listens on:    10.0.0.192:8080
TLS/SSL is            Disabled
POST API Address:     10.0.0.192:8080/nc
Websocket Address:    10.0.0.192:8080/ws
Database Path:        /tmp/data.db

```

Output above is default setting. Akari based server could detect your IP address and listen on port 8080.
> Note: In order to run the server, you need to set database path at least.

#### Trigger your custom function

```sh
curl -H "Content-Type: application/json" \
-X POST \
-d '{"Source": "f6283b29169cf8c1e84bf23cf86772fb", "Destination": ["HANDLERFUNC", "PRINT"], "Data": {"customFunc": "This is custom function content."}}' \
http://10.0.0.192:8080/nc
```

If you get `{"Status":"ok!"}`, you will see `This is custom function content.` is printed in terminal:

```sh
$ go run akariServer.go

Server listens on:    10.0.0.192:8080
TLS/SSL is            Disabled
POST API Address:     10.0.0.192:8080/nc
Websocket Address:    10.0.0.192:8080/ws
Database Path:        /tmp/data.db

This is custom function content.
[GIN] 2016/12/08 - 13:22:50 | 200 |     107.267µs | 10.0.0.192 |   POST    /nc

```
> [gin-gonic/gin](https://github.com/gin-gonic/gin) is currently used for handling HTTP requests because Akari's own HTTP handler is still under developing.

## How it works

Akari message framework serves a HTTP POST API to receive HTTP requests, and serves a websocket service to communicate with any device suppoets websocket. It receives messages from HTTP POST request or websocket. Then pushing messages to target destination(s).

Akari message framework also supports to broadcast messages; send message to third party services, such as Pushbullet; or perform a custom behavior.

## Features

### Mandatory Identity Authentication

Akari message framework identifies every device by token. Each message must have correct token of `Source` and `Destination`. When a device try to connect websocket service, it needs to provide its identity in 30 second in order to register itself. Otherwise, websocket service will reject and close the connection.

However, for human readable purpose, token is stored with `name`:

```json
{
   "Name":"Akari",
   "Token":"f6283b29169cf8c1e84bf23cf86772fb"
}
```

### Unified Message Format

All messages sent to or sent by Akari message framework has an unified format. It means all messages transferred with Akari message framework must follow this Json format:

```json
{
   "Source":"example token",
   "Destination":[
      "example token or special command",
      "example token"
   ],
   "Data":{
      "example 1":"example",
      "example 2":"example"
   }
}
```

Akari message framework reads and check `Source` and `Destination` to determine where the message is from and where the message is going. `Data` is utilized for users to exchange information.

### Broadcast

If `Destination` set as `["BROADCAST"]`, Akari based server will broadcast this message to every device registered as online.

```json
{
   "Source":"example token",
   "Destination":[
      "BROADCAST"
   ],
   "Data":{
      "example 1":"example",
      "example 2":"example"
   }
}
```

### Custom Function

In order to run a custom function set in Event map, `Destination` need to be set as `["HANDLERFUNC"]` with your Event name following.

For example, if your custom function (`PRINT`) is:

```go
func main() {
	// (Rest of main)

	// Custom function
	// "PRINT" is the Event name
	c.Event["PRINT"] = eventPrint

	// (Rest of main)
}

// Custom function
func eventPrint(m *akari.Message) error {
	fmt.Println(m.Data["customFunc"])
	return nil
}
```

Your message should look like this:

```json
{
   "Source":"example token",
   "Destination":[
      "HANDLERFUNC",
      "PRINT"
   ],
   "Data":{
      "example 1":"example",
      "example 2":"example"
   }
}
```

### Pushbullet Support

Akari message framework supports sending notification via Pushbullet. Set `"Destination":["PUSHBULLET"]` to send a message to Pushbullet service.
> Currently, only support sending "push" notification in type "note".

```json
{
   "Source":"example token",
   "Destination":[
      "PUSHBULLET",
      "example token",
      "example token"
   ],
   "Data":{
      "Type":"note",
      "Title":"push a note",
      "Body":"note body",
      "AccessToken":"your Pushbullet token"
   }
}
```
If you set multiple destinations, Akari based server will try to send the message to the destinations following `"HPUSHBULLET"`. If one of those destinations is offline, Akari based server will send the message to Pushbullet. This method could be seen as adding an alternative destination for receiving notification.
> Note: `"Data"` should have same format as example above. Otherwise, Pushbullet notification would fail to send.
