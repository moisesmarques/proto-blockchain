package grpcBlockchain

import (
	"fmt"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

func StartServer(port string) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	// Check for errors
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Instanciate the server
	s := grpc.NewServer()
	RegisterGrpcBlockchainServiceServer(s, &Server{})
	fmt.Println("GRPC SERVICE RUNNING ON PORT:", port)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("[x] serve: %v", err)
	}
}
