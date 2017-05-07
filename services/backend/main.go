package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/tcnksm/go-distributed-trace/proto/message"

	"cloud.google.com/go/trace"
	"google.golang.org/grpc"
)

const (
	EnvProjectID              = "PROJECT_ID"
	EnvMessageServiceHost     = "MESSAGE_SERVICE_HOST"
	EnvGRPCMessageServiceHost = "GRPC_MESSAGE_SERVICE_HOST"
)

func main() {

	traceClient, err := trace.NewClient(context.TODO(), os.Getenv(EnvProjectID))
	if err != nil {
		log.Fatalf("[ERROR] Failed to create new trace client: %s", err)
	}

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s %s", r.Method, r.URL.Path, r.UserAgent())
		span := traceClient.SpanFromRequest(r)
		defer span.Finish()

		heavyF := func(name string, t time.Duration) {
			spanC := span.NewChild(name)
			defer spanC.Finish()
			time.Sleep(t)
		}

		// TODO(tcnksm): Does not work .....
		heavyF("backend-process1", 10*time.Millisecond)
		heavyF("backend-process2", 30*time.Millisecond)

		host := os.Getenv(EnvMessageServiceHost)
		grpcHost := os.Getenv(EnvGRPCMessageServiceHost)

		if len(host) != 0 {
			// Do normal http request
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

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, resp.Body)

		} else if len(grpcHost) != 0 {
			// Do gRPC request
			conn, err := grpc.Dial(grpcHost, grpc.WithInsecure(), grpc.WithUnaryInterceptor(trace.GRPCClientInterceptor()))
			if err != nil {
				log.Fatalf("[ERROR] Faield to connect: %v", err)
				return
			}
			defer conn.Close()

			client := pb.NewMessageClient(conn)

			grpcSpan := span.NewChild(grpcHost)
			defer grpcSpan.Finish()

			ctx := trace.NewContext(r.Context(), grpcSpan)
			resp, err := client.Hello(ctx, &pb.HelloRequest{
				Name: "tcnksm",
			})
			if err != nil {
				log.Fatalf("[ERROR] Faield: %v", err)
				return
			}

			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp.Message + "\n"))
			return

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No backend message service is working\n"))
			return
		}
	})

	log.Printf("[INFO] Start server on :3001")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %s", err)
	}
}
