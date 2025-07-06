package application

import "WebSocket_Front/src/domain"

type WebsocketService interface {
	SendMessage(message domain.Message) error
}
