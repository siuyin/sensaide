package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"image_capture/internal/public"

	"cloud.google.com/go/vertexai/genai"
	"github.com/siuyin/dflt"
)

const projectID = "lsy0318"
const location = "us-west1"

func main() {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! It is %v\n", time.Now().Format("15:04:05.000 MST"))
	})

	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		getFile(r)

		ctx := context.Background()

		client, err := genai.NewClient(ctx, projectID, location)
		if err != nil {
			log.Println("failed to init new client: ", err)
			return
		}

		model := client.GenerativeModel("gemini-1.0-pro-vision-001")
		imgData, _ := os.ReadFile("image.jpg")
		img := genai.ImageData("jpeg", imgData)
		prompt := genai.Text("Tell me about the people, if any, in the image at Location ID:444 do they look like they can benefit from better cooling or heating. Explain your reasoning. Finally mention if additional cooling or warming is warranted by stating action: increase cooling or action: increase warming in your response. or action: no action required Also include the Location ID in your response.")
		resp, err := model.GenerateContent(ctx, img, prompt)
		if err != nil {
			log.Printf("error on generate content: %v", err)
		}
		result := respText(resp)
		log.Println(result)
		fmt.Fprintf(w, "%s", result)
	})

	// http.Handle("/", http.FileServer(http.Dir("./internal/public"))) // uncomment for development
	http.Handle("/", http.FileServer(http.FS(public.Content))) // uncomment for deployment

	log.Fatal(http.ListenAndServe(":"+dflt.EnvString("HTTP_PORT", "8080"), nil))
}

func getFile(r *http.Request) {
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

	os.WriteFile("image.jpg", b, 0666)
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
