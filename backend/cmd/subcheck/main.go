package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/siuyin/sensaide/backend/internal/msg"
)

func main() {
	ctx := context.Background()
	for {
		if err := msg.Sub.Receive(ctx, subHandler); err != nil {
			log.Printf("msg.Sub.Receive: %s", err)
		}
	}
}

func subHandler(ctx context.Context, m *pubsub.Message) {
	fmt.Printf("msg: %s\n", m.Data)
	m.Ack()
}
