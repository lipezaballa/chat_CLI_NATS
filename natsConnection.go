package main

import (
	"chat_CLI_NATS/data"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func InitJetStream(client *data.ChatClient, channels []string) (*data.ChatClient, error)  {
	log.Println("iniciar JetStream")

	js, err := client.Nc.JetStream()
	if err != nil {
		log.Fatalf("Error al inicializar JetStream: %v", err)
		return nil, err
	}
	client.Js = js

	streamName := "CHAT"
	_, err = js.StreamInfo(streamName)
	if err != nil {
		log.Printf("El stream '%s' no existe. Creando el stream...", streamName)

		// Configuración para crear el stream (puedes ajustarlo según tus necesidades)
		streamConfig := &nats.StreamConfig{
			Name:     streamName,
			Subjects: channels,
			Retention: nats.LimitsPolicy,    
			MaxMsgs:   -1,                      
			MaxBytes:  -1,                      
			MaxAge:    1 * time.Hour,                       
			Storage:   nats.FileStorage,      
		}

		// Intentamos crear el stream
		_, err = js.AddStream(streamConfig)
		if err != nil {
			log.Fatalf("Error creando el stream: %v", err)
		}

		log.Printf("Stream '%s' con canales '%s' creado exitosamente", streamName, channels)
	} else {
		log.Printf("El stream '%s' ya existe", streamName)
	}

	return client, nil
}