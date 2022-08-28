package grpcBlockchain

import (
	"encoding/json"
	"flag"
	"fmt"
	"judge-blockchain/blockchain"
	network "judge-blockchain/communication/tcp"
	"judge-blockchain/wallet"
	"log"
	"strconv"
	"time"

	grpcAuditLog "judge-blockchain/communication/grpc/schema/auditLog"
	pb "judge-blockchain/communication/grpc/schema/cpu"

	consensusPkg "judge-blockchain/consensus"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
)

type Server struct {
	UnimplementedGrpcBlockchainServiceServer
}

func (s *Server) PrintChain(ctx context.Context, in *EmptyRequest) (*SuccessResponse, error) {
	response := ""
	log.Println("GRPC WORKING")
	if !blockchain.DBexists() {
		return &SuccessResponse{Body: "Blockchain does not exist"}, nil
		// runtime.Goexit()
	} else {

		chain := blockchain.ContinueBlockChain()
		defer chain.Database.Close()
		iterator := chain.Iterator()
		for {
			block := iterator.Next()
			response += "Previous hash: " + fmt.Sprint(block.PrevHash) + "\n"
			response += "hash: " + fmt.Sprint(block.Hash) + "\n"
			// fmt.Printf("hash: %x\n", block.Hash)
			pow := blockchain.NewProofOfWork(block)
			response += "Pow: " + strconv.FormatBool(pow.Validate()) + "\n"
			for _, tx := range block.Transactions {
				fmt.Println(tx)
				response += "Transaction: " + fmt.Sprint(tx.ID) + "\n"
			}
			response += "Height: " + strconv.Itoa(block.Height) + "\n"
			response += "TimeStamp: " + fmt.Sprint(block.Timestamp) + "\n"
			response += "Nonce: " + strconv.Itoa(block.Nonce) + "\n"
			response += "*************************************************\n"
			response += "*******************************************************************************\n"
			if len(block.PrevHash) == 0 {
				break
			}

		}
	}
	return &SuccessResponse{Body: response}, nil
}

func (s *Server) CreateChain(ctx context.Context, in *Wallet) (*SuccessResponse, error) {
	address := in.Body
	if !wallet.ValidateAddress(address) {
		return &SuccessResponse{Body: "Invalid Address"}, nil
	}
	if blockchain.DBexists() {
		return &SuccessResponse{Body: "Blockchain already exists"}, nil
	}
	newChain := blockchain.InitBlockChain(address)
	newChain.Database.Close()
	// fmt.Println("Finished creating chain")
	return &SuccessResponse{Body: "Finished creating chain"}, nil
}

func (s *Server) CreateActionData(ctx context.Context, in *ActionDataRequest) (*SuccessResponse, error) {
	toJson := protojson.Format(in)
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	var actionResult blockchain.ValidatorActionResult
	json.Unmarshal([]byte(toJson), &actionResult)
	tx := blockchain.NewActionResult(chain, actionResult)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
	response := ""
	iterator := chain.Iterator()
	// for {
	// 	block := iterator.Next()
	// 	response += "Previous hash: " + string(block.PrevHash) + "\n"
	// 	response += "hash: " + string(block.Hash) + "\n"
	// 	pow := blockchain.NewProofOfWork(block)
	// 	response += "Pow: " + strconv.FormatBool(pow.Validate()) + "\n"
	// 	fmt.Println()
	// 	if len(block.PrevHash) == 0 {
	// 		break
	// 	}
	// }
	lastBlock := iterator.Next()
	response += "Previous hash: " + fmt.Sprint(lastBlock.PrevHash) + "\n"
	response += "hash: " + fmt.Sprint(lastBlock.Hash) + "\n"
	var blocks []blockchain.Block
	blocks = append(blocks, *lastBlock)
	for _, peer := range network.KnownNodes {
		network.SendBlocks(peer, "blocks", blocks)
	}
	return &SuccessResponse{Body: response}, nil

}

// function added by Sameer for handling reception of newly generated shell schema

func (*Server) HandleNewShellSchema(ctx context.Context, shell *NewShellSchema) (*NewShellSchemaResponse, error) {
	log.Printf("Received Shell Schema: %s", shell)
	tmp := NewShellSchemaResponse{JsonFAK: "{test:\"FAK1\"}"}

	return &tmp, nil
}

func (s *Server) ProcessAdminRequest(ctx context.Context, in *AdminRequest) (*AdminReply, error) {
	log.Printf("FAK Received: %v", in.GetFak())

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCPUShardClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ProcessRequest(ctx, &pb.CPURequest{ActionId: "34b631b1-a8b1-49cd-937a-2d7185909e28",
		ActionName:   "gpu_task",
		Function:     "XXXXX",
		SecreteShare: "XXXYYYYYY",
		Sender:       "XXXYYYYYY",
		ParticipantingNodes: []string{"34b631b1-a8b1-49cd-937a-2d7185909e28",
			"34b631b1-a8b1-49cd-937a-2d7185909e29",
			"34b631b1-a8b1-49cd-937a-2d7185909e30",
			"34b631b1-a8b1-49cd-937a-2d7185909e31",
			"34b631b1-a8b1-49cd-937a-2d7185909e32",
			"34b631b1-a8b1-49cd-937a-2d7185909e33",
		},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetActionName())

	return &AdminReply{Message: "Respnse for - " + in.GetFak()}, nil

}

func (s *Server) SendAuditLogToAdmin(ctx context.Context, data *grpcAuditLog.SendAuditToAdminRequest) (*grpcAuditLog.SendAuditToAdminResponse, error) {

	fmt.Printf("Received Audit summary from storage: %v", data)

	consensusData, _ := consensusPkg.CreateConsensusData()

	auditLog := consensusPkg.CreateAuditLog(data.Auditlog.Identity,
		data.Auditlog.Chunk,
		int32(data.Auditlog.TimeInMilliSeconds),
		data.Auditlog.Result,
		data.Auditlog.Message,
		data.Auditlog.NextAction,
	)

	consensusData.AddAuditLog(auditLog)

	fmt.Printf("consensus auditLog %v \n", auditLog)
	consensusData.SaveFile()

	return &grpcAuditLog.SendAuditToAdminResponse{Message: "OK"}, nil
}

func (s *Server) SendAuditToAdmin(ctx context.Context, data *grpcAuditLog.SendAuditToAdminRequest) (*grpcAuditLog.SendAuditToAdminResponse, error) {

	fmt.Printf("Received Audit log from storage: %v", data)

	consensusData, _ := consensusPkg.CreateConsensusData()

	audit := consensusPkg.CreateAudit(data.Auditsummary.Identity,
		int32(data.Auditsummary.Successes),
		int32(data.Auditsummary.Offlines),
		int32(data.Auditsummary.Fails),
		data.Auditsummary.Status,
		int32(data.Auditsummary.TotalCount),
		fmt.Sprintf("%v", data.Auditsummary.LastAuditTime),
	)

	consensusData.AddAudit(audit)

	fmt.Printf("consensus audit %v \n", audit)
	consensusData.SaveFile()

	return &grpcAuditLog.SendAuditToAdminResponse{Message: "OK"}, nil

}
