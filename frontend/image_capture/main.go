package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"image_capture/internal/public"

	"cloud.google.com/go/vertexai/genai"
	"github.com/siuyin/dflt"
)

const projectID = "lsy0318"
const location = "us-west1"

type vertexResponse struct {
	Action        string
	LocationID    string
	Reason        string
	PeoplePresent bool
}

type controlRoomMessage struct {
	RoomID        string
	Text          string `json:"text"`
	Device        string
	OnOff         int `json:"on_off"`
	PeoplePresent bool
}

func main() {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! It is %v\n", time.Now().Format("15:04:05.000 MST"))
	})

	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		roomID := r.FormValue("roomID")

		b := getBytes(r)
		res := callVertexAI(b, roomID)

		fmt.Fprintf(w, "%s", res)
	})

	// http.Handle("/", http.FileServer(http.Dir("./internal/public"))) // uncomment for development
	t := template.Must(template.ParseFS(public.Content, "index.html"))
	http.HandleFunc("/roomid", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		t.Execute(w, struct{ RoomID string }{id})
	})
	http.Handle("/", http.FileServer(http.FS(public.Content))) // uncomment for deployment

	log.Fatal(http.ListenAndServe(":"+dflt.EnvString("HTTP_PORT", "8080"), nil))
}

func callVertexAI(b []byte, roomID string) string {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Println("failed to init new client: ", err)
		return ""
	}

	model := client.GenerativeModel("gemini-1.0-pro-vision-001")
	model.SetTemperature(0.3)
	img := genai.ImageData("jpeg", b)
	prompt := genai.Text(fmt.Sprintf(`Tell me about the people, if any (PeoplePresent true or false), in the image at Location ID:"%s" do they look like they can benefit from better cooling or heating. Explain your reasoning. Finally mention if additional cooling or warming is warranted by stating action: increase cooling or action: increase warming in your response. or action: no action required Also include the Location ID in your response. output in JSON format with fields PeoplePresent, Action, Reason, LocationID`, roomID))
	resp, err := model.GenerateContent(ctx, img, prompt)
	if err != nil {
		log.Printf("error on generate content: %v", err)
	}
	result := respText(resp)

	result = stripHeaderAndFooter(result)

	msg := vertexResponse{}
	if err := json.Unmarshal([]byte(result), &msg); err != nil {
		log.Printf("failed to unmarshal result: %s", err)
	}

	recommend(&msg, roomID)

	return result
}

func recommend(msg *vertexResponse, roomID string) {
	crMsg := controlRoomMessage{}
	crMsg.Device = "aircon"
	crMsg.RoomID = roomID
	crMsg.Text = msg.Reason

	if !msg.PeoplePresent {
		crMsg.OnOff = 0
		updateControlRoom(&crMsg)
		crMsg.Device = "lights"
		updateControlRoom(&crMsg)

		return
	}

	if msg.PeoplePresent {
		crMsg.OnOff = 1
		crMsg.Device = "lights"
		updateControlRoom(&crMsg)
	}

	if msg.Action == "increase heating" {
		crMsg.OnOff = 0
		updateControlRoom(&crMsg)
	} else {
		crMsg.OnOff = 1
		updateControlRoom(&crMsg)
	}
}

func getBytes(r *http.Request) []byte {
	r.ParseMultipartForm(1000000)
	fh := r.MultipartForm.File["data"]

	f, err := fh[0].Open()
	if err != nil {
		log.Printf("file header open failed: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Printf("io read all failed: %v", err)
	}

	return b
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

func stripHeaderAndFooter(s string) string {
	r := strings.ReplaceAll(s, "```json", "")
	r = strings.ReplaceAll(r, "```", "")

	return r
}

func updateControlRoom(msg *controlRoomMessage) error {
	URL := fmt.Sprintf("https://spgroup24.alwaysdata.net/%v/%v", msg.Device, msg.RoomID)
	log.Println(URL)

	dat, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("updateControlRoom: json.Marshal: %s", err)
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(dat))
	if err != nil {
		return fmt.Errorf("updateControlRoom: http.Post: %v: %v", resp, err)
	}

	return nil
}
