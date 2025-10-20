package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"rpc/internal/server"
	"rpc/pkg/api/test"
)

func main() {
	serv := grpc.NewServer()
	orderServer := server.NewServer()
	test.RegisterOrderServiceServer(serv, orderServer)
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on %v:%v\n", l.Addr(), err)
	}
	log.Printf("Trying to start grpc order server on %v\n", l.Addr())
	if err := serv.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
