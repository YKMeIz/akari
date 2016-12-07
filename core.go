package akari

import (
	"github.com/gin-gonic/gin"
)

type Core struct {
	Domain                string
	Port                  string
	CertChain             string
	CertKey               string
	MessageRelativePath   string
	WebsocketRelativePath string
	DatabasePath          string
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
	r.GET(c.WebsocketRelativePath, func(c *gin.Context) {
		handleWebsocket(hub, c.Writer, c.Request)
	})
	r.POST(c.MessageRelativePath, func(g *gin.Context) {
		handleMsg(hub, g)
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
