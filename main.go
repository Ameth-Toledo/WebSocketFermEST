package main

import (
	"WebSocket_Front/src/infraestructure/mqtt"
	"WebSocket_Front/src/infraestructure/routes"
	"WebSocket_Front/src/infraestructure/webSocket"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	mqttSender := mqtt.NewMQTTSender()
	wsService := webSocket.NewWebsocketService(mqttSender)

	router := gin.Default()

	routes.RegisterRoutes(router, wsService)

	log.Println("Servidor WebSocket escuchando en :8088")
	if err := router.Run(":8088"); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
