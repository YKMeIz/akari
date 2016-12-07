package akari

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	// Message sender's token.
	Source string

	// Message receiver's token.
	Destination []string

	// Message content.
	Data map[string]string
}

func handleMsg(hub *Hub, g *gin.Context) {
	if err := sendToWebsocket(g, hub); err != nil {
		g.JSON(http.StatusNotImplemented, gin.H{"Status": "error! " + err.Error()})
		g.AbortWithStatus(http.StatusNotImplemented)
		return
	}
	g.JSON(http.StatusOK, gin.H{"Status": "ok!"})
}
