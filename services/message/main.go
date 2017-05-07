package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/trace"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())
		span := trace.FromContext(r.Context())
		defer span.Finish()

		time.Sleep(10 * time.Millisecond)

		response := struct {
			Message string `json:"message"`
		}{
			Message: "hello, world",
		}

		w.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(&response); err != nil {
			log.Printf("[ERROR] Failed to encode message: %s", err)
			return
		}
	})

	log.Printf("[INFO] Start server")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server")
	}
}
