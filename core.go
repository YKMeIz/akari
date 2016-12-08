package akari

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
)

func (m *Message) runHandlerFunc(f HandlerFunc) error {
	return f(m)
}

func New() *Core {
	ip, err := IPAddress()
	if err != nil {
		panic(err.Error())
	}
	c := &Core{
		Domain:                ip,
		Port:                  "8080",
		MessageRelativePath:   "/nc",
		WebsocketRelativePath: "/ws"}
	c.Event = make(map[string]HandlerFunc)
	return c
}

func (c Core) Run() {
	c.isDatabasePath()
	c.OpenDatabase()
	c.printInfo()
	c.serve()
	defer db.Close()
}

func (c Core) isDatabasePath() {
	if c.DatabasePath == "" {
		panic("err: Database path is not set.")
	}
}

func (c Core) printInfo() {
	var https string
	if c.CertChain != "" && c.CertKey != "" {
		https = "\033[32m\033[1menabled\033[0m\033[39m"
	} else {
		https = "\033[31m\033[1mDisabled\033[0m\033[39m"
	}
	fmt.Println("\nServer listens on:    " + c.Domain + ":" + c.Port)
	fmt.Println("TLS/SSL is            " + https)
	fmt.Println("POST API Address:     " + c.Domain + ":" + c.Port + c.MessageRelativePath)
	fmt.Println("Websocket Address:    " + c.Domain + ":" + c.Port + c.WebsocketRelativePath)
	fmt.Println("Database Path:        " + c.DatabasePath + "\n")
}

func (c Core) serve() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	hub := newHub()
	go hub.run()
	r.GET(c.WebsocketRelativePath, func(g *gin.Context) {
		c.handleWebsocket(hub, g.Writer, g.Request)
	})
	r.POST(c.MessageRelativePath, func(g *gin.Context) {
		c.handleApi(hub, g)
	})
	if c.CertChain == "" || c.CertKey == "" {
		if c.Port == "" {
			r.Run(c.Domain + ":80")
		} else {
			r.Run(c.Domain + ":" + c.Port)
		}
	} else {
		if c.Port == "" {
			r.RunTLS(c.Domain+":443", c.CertChain, c.CertKey)
		} else {
			r.RunTLS(c.Domain+":"+c.Port, c.CertChain, c.CertKey)
		}
	}
}

// handleWebsocket handles websocket requests from the peer.
func (c Core) handleWebsocket(h *hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go c.tokenRegister(h, w, r, conn)
}

func (c Core) handleApi(h *hub, g *gin.Context) {
	if err := c.sendToWebsocket(g, h); err != nil {
		g.JSON(http.StatusNotImplemented, gin.H{"Status": "error! " + err.Error()})
		g.AbortWithStatus(http.StatusNotImplemented)
		return
	}
	g.JSON(http.StatusOK, gin.H{"Status": "ok!"})
}

func IPAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// interface down
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			// loopback interface
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// not an ipv4 address
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("err: No network connection.")
}
