docker run --name nats -it -p 4222:4222 nats --js

go run main.go nats://localhost:4222 chat Felipe
