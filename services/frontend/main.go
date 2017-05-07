package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/trace"
)

const (
	EnvProjectID   = "PROJECT_ID"
	EnvBackendHost = "BACKEND_HOST"
)

type messageClient struct{}

func main() {

	traceClient, err := trace.NewClient(context.TODO(), os.Getenv(EnvProjectID))
	if err != nil {
		log.Fatalf("[ERROR] Failed to create new trace client: %s", err)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !(username == "tcnksm" && password == "ncpai4wbp948bc49qu") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())
		span := traceClient.NewSpan("/hello")
		defer span.Finish()

		heavyF := func(name string, t time.Duration) {
			spanC := span.NewChild(name)
			defer spanC.Finish()
			time.Sleep(t)
		}

		heavyF("frontend-process1", 40*time.Millisecond)
		heavyF("frontend-process2", 50*time.Millisecond)

		host := os.Getenv(EnvBackendHost)
		if len(host) == 0 {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No backend service is working\n"))
			return
		}

		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to create new request: %s", err)
			return
		}

		remoteSpan := span.NewRemoteChild(req)
		defer remoteSpan.Finish()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[ERROR] Failed to do request: %s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, resp.Body)
	})

	log.Printf("[INFO] Start server")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %s", err)
	}
}
