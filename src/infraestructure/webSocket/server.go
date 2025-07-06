package webSocket

import (
	"WebSocket_Front/src/application"
	"WebSocket_Front/src/domain"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type WebsocketService struct {
	sender      application.WebsocketService
	connections map[int][]*websocket.Conn
	mu          sync.Mutex
}

func NewWebsocketService(sender application.WebsocketService) *WebsocketService {
	return &WebsocketService{
		sender:      sender,
		connections: make(map[int][]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WebsocketService) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al establecer conexión WebSocket: %v", err)
		return
	}

	// Obtener user_id desde query
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		log.Println("user_id faltante en la URL")
		conn.Close()
		return
	}
	var userID int
	fmt.Sscanf(userIDStr, "%d", &userID)

	// Registrar la conexión
	ws.mu.Lock()
	ws.connections[userID] = append(ws.connections[userID], conn)
	ws.mu.Unlock()
	log.Printf("✅ Conectado user_id=%d. Conexiones actuales: %d", userID, len(ws.connections[userID]))

	defer func() {
		ws.removeConnection(userID, conn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("❌ Error leyendo mensaje WebSocket: %v", err)
			break
		}

		var message domain.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("⚠️  Mensaje inválido recibido: %v", err)
			continue
		}

		log.Printf("📨 Mensaje recibido de WebSocket: %+v", message)

		err = ws.sender.SendMessage(message)

		var response map[string]string
		if err != nil {
			response = map[string]string{
				"status":  "error",
				"message": "Error al enviar a MQTT",
			}
		} else {
			response = map[string]string{
				"status":  "ok",
				"message": "Mensaje enviado a MQTT",
			}
		}

		respJSON, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, respJSON)

		// Enviar a todos los conectados con el mismo user_id
		ws.BroadcastToUser(userID, message)
	}
}

func (ws *WebsocketService) removeConnection(userID int, conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	conns := ws.connections[userID]
	for i, c := range conns {
		if c == conn {
			ws.connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	log.Printf("🔌 Conexión cerrada user_id=%d. Conexiones restantes: %d", userID, len(ws.connections[userID]))
}

func (ws *WebsocketService) BroadcastToUser(userID int, message domain.Message) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	conns := ws.connections[userID]
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Error serializando para broadcast: %v", err)
		return
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("❌ Error enviando mensaje a conexión: %v", err)
		}
	}
}
