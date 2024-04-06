package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/vertexai/genai"
	"github.com/siuyin/dflt"
)

func main() {
	http.HandleFunc("/api", msgHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type sensorMsg struct {
	RoomID    string
	Timestamp string
	URL       []string
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	msg := sensorMsg{}
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("post read: %v", err)
	}

	if err := json.Unmarshal(dat, &msg); err != nil {
		fmt.Fprintf(w, "json.Unmarshal: %v", err)
		return
	}

	log.Println(msg)
	resp, err := callVertexAI(&msg)
	if err != nil {
		fmt.Fprintf(w, "callVertexAI: %v", err)
	}
	updControlRoom(resp)
}

type controlRoomMsg struct {
	RoomID string
	Text   string `json:"text"`
	Device string
	OnOff  int `json:"on_off"`
}

func callVertexAI(msg *sensorMsg) (*controlRoomMsg, error) {
	var projectID = dflt.EnvString("PROJECT_ID", "myproject")
	var location = dflt.EnvString("VERTEX_LOCATION", "us-west1")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return &controlRoomMsg{}, fmt.Errorf("genai.NewClient: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.0-pro-vision-001")
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text("what is the meaning of life?"))
	if err != nil {
		return &controlRoomMsg{}, fmt.Errorf("model.GenerateContent: %v", err)
	}

	return &controlRoomMsg{
		RoomID: msg.RoomID,
		Text:   respText(resp),
		OnOff:  1,
		Device: "aircon",
	}, nil
}
func respText(resp *genai.GenerateContentResponse) string {
	s := ""
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			s = s + fmt.Sprintf("%s", part)
		}
	}
	return s
}

func updControlRoom(msg *controlRoomMsg) error {
	//const URL = "https://spgroup24.alwaysdata.net/lights/34"
	URL := fmt.Sprintf("https://spgroup24.alwaysdata.net/%v/%v", msg.Device, msg.RoomID)
	log.Println(URL)
	//buf := bytes.NewBufferString(fmt.Sprintf("{\"text\":\"%s\",\"on_off\":%d}", msg.Text, msg.OnOff))
	dat, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("updateControlRoom: json.Marshal: %", err)
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(dat))
	if err != nil {
		return fmt.Errorf("updateControlRoom: http.Post: %v: %v", resp, err)
	}

	return nil
}
