package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/siuyin/sensaide/backend/internal/msg"
)

func main() {
	msg.Sub.Receive(context.Background(), msgHandler)
}

func msgHandler(ctx context.Context, m *pubsub.Message) {
	// handle incoming messages
	msg := parse(m)
	for _, u := range msg.URL {
		fmt.Println(u)
	}
	// call Vertex AI
	// update Control Room UI
	updControlRoom()
}

type pubsubMsg struct {
	RoomID    string
	Timestamp string
	URL       []string
}

func parse(m *pubsub.Message) *pubsubMsg {
	msg := pubsubMsg{}
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		log.Printf("parse: %v", err)
	}

	m.Ack()
	return &msg
}

func updControlRoom() {
	const URL = "https://spgroup24.alwaysdata.net/lights/34"
	buf := bytes.NewBufferString("{\"text\":\"test message\",\"on_off\":1}")
	http.Post(URL, "application/json", buf)
}
