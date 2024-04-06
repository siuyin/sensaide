package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/siuyin/dflt"
	"image_capture/internal/public"
)

func main() {
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! It is %v\n", time.Now().Format("15:04:05.000 MST"))
	})

	// http.Handle("/", http.FileServer(http.Dir("./internal/public"))) // uncomment for development
	http.Handle("/", http.FileServer(http.FS(public.Content))) // uncomment for deployment

	log.Fatal(http.ListenAndServe(":"+dflt.EnvString("HTTP_PORT", "8080"), nil))
}
