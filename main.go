package main

import (
	"bufio"
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
	// Chat channel
	channel := os.Args[2]
	// User name
	name := os.Args[3]

	// Connect to NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error al conectar con NATS: %v", err)
	}
	defer nc.Close()

	// Subscribe to the channel
	_, err = nc.Subscribe(channel, func(msg *nats.Msg) {
		// Show received messages
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Error al suscribirse al canal: %v", err)
	}

	fmt.Printf("Conectado al chat en el canal '%s'.\n", channel)

	// Read messages written by the user on the terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		//Problem: Si estoy escribiendo y llega un mensaje lo escribe en mi línea, luego cuando envías se envia todo bien, pero quedar raro
		if strings.TrimSpace(text) == "exit" {
			fmt.Println("Saliendo del chat...")
			break
		}
		// Publish the message in the channel
		message := fmt.Sprintf("[%s]: %s", name, text)
		if err := nc.Publish(channel, []byte(message)); err != nil {
			log.Printf("Error al enviar el mensaje: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error al leer entrada: %v", err)
	}
}