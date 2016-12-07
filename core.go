package akari

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func runHandlerFunc(f HandlerFunc) error {
	return f()
}

func (c *Core) InitEventHandler() {
	c.Event = make(map[string]HandlerFunc)
}

func (c Core) Run() {
	initDatabase(c.DatabasePath)
	c.checkNecessaryVariable()
	c.serve()
	defer db.Close()
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

func (c Core) checkNecessaryVariable() {
	if c.Domain == "" {
		panic("err: Domain is missing.")
	}
	if c.MessageRelativePath == "" {
		panic("err: relative path of message api is missing.")
	}
	if c.WebsocketRelativePath == "" {
		panic("err: relative path of websocket is missing.")
	}
}

// handleWebsocket handles websocket requests from the peer.
func (c Core) handleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go c.tokenRegister(hub, w, r, conn)
}

func (c Core) handleApi(hub *Hub, g *gin.Context) {
	if err := c.sendToWebsocket(g, hub); err != nil {
		g.JSON(http.StatusNotImplemented, gin.H{"Status": "error! " + err.Error()})
		g.AbortWithStatus(http.StatusNotImplemented)
		return
	}
	g.JSON(http.StatusOK, gin.H{"Status": "ok!"})
}
