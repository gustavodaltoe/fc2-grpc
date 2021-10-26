package main

import (
	"log"
	"net"

	"github.com/gustavodaltoe/fc2-grpc/pb"
	"github.com/gustavodaltoe/fc2-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051") // listener
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer()                                      // server
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService()) // register service
	reflection.Register(grpcServer)                                     // register reflection

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}
}
