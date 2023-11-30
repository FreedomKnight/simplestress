package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/FreedomKnight/simplestress/proto"
	"google.golang.org/grpc"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

type server struct {
    pb.UnimplementedPaddleServer
}

func (s *server) Serve(c context.Context, ping *pb.Ping) (*pb.Pong, error) {
    log.Printf("Received Ping %s\n", ping.GetMessage())
    return &pb.Pong{Message: ping.GetMessage()}, nil
}


func main() {
    flag.Parse()
    lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    s := grpc.NewServer()
    pb.RegisterPaddleServer(s, &server{})
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
