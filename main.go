package main

import (
	"bufio"
	"chat_CLI_NATS/data"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

func main() {
	// Validate arguments
	if len(os.Args) != 4 {
		log.Fatalf("Uso: %s <nats-url> <canal> <nombre>", os.Args[0])
	}

	// NATS server IP
	natsURL := os.Args[1]

	// Create a ChatClient instance
    client := &data.ChatClient{
		// Chat channel
        Channel: os.Args[2],
		// User name
        Name:    os.Args[3],
    }

	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error al conectar con NATS: %v", err)
	}
	defer nc.Close()

	client.Nc = nc

	// Subscribe to the channel
	_, err = nc.Subscribe(client.Channel, func(msg *nats.Msg) {
		// Show received messages
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Error al suscribirse al canal: %v", err)
	}

	fmt.Printf("Conectado al chat en el canal '%s'.\n", client.Channel)

	// Read messages written by the user on the terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		//Problem: Si estoy escribiendo y llega un mensaje lo escribe en mi línea, luego cuando envías se envia todo bien, pero quedar raro
		if strings.TrimSpace(text) == "exit" {
			exitChat(client)
			break
		}
		// Publish the message in the channel
		publishMessage(client, text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error al leer entrada: %v", err)
	}
}

func publishMessage(client *data.ChatClient, text string) {
	message := fmt.Sprintf("[%s]: %s", client.Name, text)
		if err := client.Nc.Publish(client.Channel, []byte(message)); err != nil {
			log.Printf("Error al enviar el mensaje: %v", err)
		}
}

func exitChat(client *data.ChatClient) {
	message := fmt.Sprintf("%s salió del chat...\n", client.Name)
	if err := client.Nc.Publish(client.Channel, []byte(message)); err != nil {
		log.Printf("Error al enviar el mensaje: %v", err)
	}
}