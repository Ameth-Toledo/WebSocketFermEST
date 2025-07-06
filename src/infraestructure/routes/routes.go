package routes

import (
	"WebSocket_Front/src/infraestructure/webSocket"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, ws *webSocket.WebsocketService) {
	router.GET("/ws", func(c *gin.Context) {
		ws.HandleConnection(c)
	})
}
