package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/tcnksm/go-distributed-trace/proto/message"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	EnvProjectID          = "PROJECT_ID"
	EnvMessageServiceHost = "MESSAGE_SERVICE_HOST"
)

type messageServer struct {
}

func (s *messageServer) Hello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	name := req.Name
	log.Printf("[INFO] Request: %s", name)
	return &pb.HelloResponse{
		Message: fmt.Sprintf("hello, %s (gRPC)", name),
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterMessageServer(grpcServer, &messageServer{})

	log.Println("[INFO] Listening grpc server on :4002")
	listener, err := net.Listen("tcp", ":4002")
	if err != nil {
		log.Fatalf("[ERROR] Failed to listen: %s", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("[ERROR] Failed to start server")
	}
}
