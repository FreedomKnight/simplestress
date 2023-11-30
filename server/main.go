package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/freedomknight/simplestress/proto"
	"google.golang.org/grpc"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

type server struct {
    pb.UnimplementedPaddleServer
}

func (s *server) Serve() (*pb.Pong, error) {
    log.Println("Received Ping")
    return &pb.Pong{message: "pong"}, nil
}


func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    s := grpc.NewServer()
    pb.RegisterPaddleServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
