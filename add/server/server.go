package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hsmtkk/cuddly-train/add/proto"
	"github.com/hsmtkk/cuddly-train/env"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedAddSerivceServer
}

func main() {
	port, err := env.Port()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("net.Listen failed; %v", err.Error())
	}
	s := grpc.NewServer()
	proto.RegisterAddSerivceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf(".Serve failed; %v", err.Error())
	}
}
