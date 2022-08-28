package main

import (
	"log"
	"time"

	"judge-blockchain/communication/grpc/schema/grpcBlockchain"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := grpcBlockchain.NewGrpcBlockchainServiceClient(conn)
	// .NewChatServiceClient(conn)

	response, err := c.PrintChain(context.Background(), &grpcBlockchain.EmptyRequest{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling PrintChain: %s", err)
	}
	log.Printf("Response from server")
	log.Printf("Response from server: %s", response.Body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ProcessAdminRequest(ctx, &grpcBlockchain.AdminRequest{ActionName: "Test action", Fak: "FAKtestFAK", ActionId: "actionIdX"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

}
