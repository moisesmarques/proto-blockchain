package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "judge-blockchain/communication/grpc/schema/cpu"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedCPUShardServer
}

func generateCPUConsensusResult() *pb.CPUConsensusResult {
	return &pb.CPUConsensusResult{

		ActionId:   "34b631b1-a8b1-49cd-937a-2d7185909e28",
		ActionName: "gpu_task",
		SenderNode: "QmPooRTbgc4tNfnRqDc9VZLnPcvAzBdowhgzeNMh4H2rfF",
		Verdict: &pb.CPUConsensusResult_Verdict{Result: "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
			ConsensusStatus: "Consent",
			NodeCount:       2,
			Nodes: []*pb.CPUConsensusResult_Nodes{{Id: "34b631b1-a8b1-49cd-937a-2d7185909e28",
				ResultHash:              "84b7431a9f4d0266319a42bced98acde0d28e604f3f29b84d4df20068c7b4696",
				DeterministicResultHash: "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
				Elapsed:                 "12.541479ms",
				DeterministicResultRaw:  "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
				ResultRaw:               "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
				NodeAddress:             "QmaeA1wFCDNtbswphc7X5ALwPmURsvNccP3w1huJw9RBve",
			},
				{Id: "34b631b1-a8b1-49cd-937a-2d7185909e27",
					ResultHash:              "84b7431a9f4d0266319a42bced98acde0d28e604f3f29b84d4df20068c7b4696",
					DeterministicResultHash: "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
					Elapsed:                 "112.541479ms",
					DeterministicResultRaw:  "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
					ResultRaw:               "4306f847af5158bdcd22387a598bd29ef81f097c2dedd2a427ffcbfd253a4568",
					NodeAddress:             "QmaeA1wFCDNtbswphc7X5ALwPmURsvNccP3w1huJw9RBve",
				},
			},
		},
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) ProcessRequest(ctx context.Context, in *pb.CPURequest) (*pb.CPUConsensusResult, error) {
	log.Printf("Received ActionId : %v", in.GetActionId())
	return generateCPUConsensusResult(), nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCPUShardServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
