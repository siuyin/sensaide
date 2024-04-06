package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api", msgHandler)
	// msg.Sub.Receive(context.Background(), msgHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type controlRoomMsg struct {
	RoomID string
	Text   string `json:"text"`
	Device string
	OnOff  int `json:"on_off"`
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	msg := pubsubMsg{}
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("post read: %v", err)
	}

	if err := json.Unmarshal(dat, &msg); err != nil {
		log.Printf("parse: %v", err)
	}

	resp := callVertexAI(&msg)
	updControlRoom(resp)
}

func callVertexAI(msg *pubsubMsg) *controlRoomMsg {
	return &controlRoomMsg{
		RoomID: msg.RoomID,
		Text:   "test message from vertex AI call",
		OnOff:  1,
		Device: "aircon",
	}
}

type pubsubMsg struct {
	RoomID    string
	Timestamp string
	URL       []string
}

func updControlRoom(msg *controlRoomMsg) {
	//const URL = "https://spgroup24.alwaysdata.net/lights/34"
	URL := fmt.Sprintf("https://spgroup24.alwaysdata.net/%v/%v", msg.Device, msg.RoomID)
	buf := bytes.NewBufferString(fmt.Sprintf("{\"text\":\"%s\",\"on_off\":%d}", msg.Text, msg.OnOff))
	http.Post(URL, "application/json", buf)
}
