package main

import (
	"bufio"
	"chat_CLI_NATS/data"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Validate arguments
	if len(os.Args) != 4 {
		log.Fatalf("Use: %s <nats-url> <canal> <nombre>", os.Args[0])
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
		log.Fatalf("Error connecting with NATS: %v", err)
	}
	defer nc.Close()

	client.Nc = nc

	//Configure JetStream
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error initializing JetStream: %v", err)
	}

	// Recover messages from last hour (not needed because stream only persist message from last hour, but done to be ensured in case stream persist everything)
	startTime := time.Now().Add(-1 * time.Hour)
	subOpts := []nats.SubOpt{
		nats.StartTime(startTime),
	}

	// Subscribe to the channel
	sub, err := js.Subscribe(client.Channel, func(msg *nats.Msg) {
		// Show received messages
		fmt.Println(string(msg.Data))
	}, subOpts...)
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Printf("Connecting to channel '%s'.\n", client.Channel)

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
		log.Fatalf("Error reading input: %v", err)
	}
}

func publishMessage(client *data.ChatClient, text string) {
	timestamp := time.Now().Format("02/01/2006 15:04:05")
	message := fmt.Sprintf("[%s] %s: %s", timestamp, client.Name, text)
		if err := client.Nc.Publish(client.Channel, []byte(message)); err != nil {
			log.Printf("Error sending message: %v", err)
		}
}

func exitChat(client *data.ChatClient) {
	message := fmt.Sprintf("%s left the chat...\n", client.Name)
	if err := client.Nc.Publish(client.Channel, []byte(message)); err != nil {
		log.Printf("Error sending exit message: %v", err)
	}
}