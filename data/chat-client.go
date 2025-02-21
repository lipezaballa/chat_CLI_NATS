package data

import (
	"github.com/nats-io/nats.go"
)

type ChatClient struct {
    Nc      *nats.Conn
    Channel string
    Name    string
    Js      nats.JetStreamContext
}