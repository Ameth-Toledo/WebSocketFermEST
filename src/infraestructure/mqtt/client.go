package mqtt

import (
	"WebSocket_Front/src/domain"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTSender struct {
	client mqtt.Client
	topic  string
}

func NewMQTTSender() *MQTTSender {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://52.202.0.30:1883")
	opts.SetClientID("mi-consumidor")
	opts.SetUsername("milton")
	opts.SetPassword("milton123")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error al conectar con MQTT: %v", token.Error())
	}

	log.Println("MQTT conectado correctamente")

	return &MQTTSender{
		client: client,
		topic:  "front",
	}
}

func (s *MQTTSender) SendMessage(message domain.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error serializando mensaje: %v", err)
		return err
	}

	token := s.client.Publish(s.topic, 0, false, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Error publicando en MQTT: %v", token.Error())
		return token.Error()
	}

	log.Printf("Mensaje MQTT enviado al tópico '%s': %s", s.topic, string(payload))
	return nil
}
