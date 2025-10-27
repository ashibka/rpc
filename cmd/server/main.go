package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/pkg/api/test"
)

func main() {
	cfg, err := config.ParseConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	grpcserver := grpc.NewServer()
	orderServer := server.NewServer()
	test.RegisterOrderServiceServer(grpcserver, orderServer)
	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen on %v:%v\n", listener.Addr(), err)
	}
	log.Printf("Trying to start grpc order server on %v\n", listener.Addr())
	if err := grpcserver.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
