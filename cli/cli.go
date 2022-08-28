package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"judge-blockchain/communication/grpc/schema/grpcBlockchain"
	api "judge-blockchain/communication/http"

	network "judge-blockchain/communication/tcp"
)

type CommandLine struct{}

var port string

//printUsage will display what options are availble to the user
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("go run main.go startAsTCP --port {{portNumber}} Start as TCP server listening on port 'portNumber'")
	fmt.Println("go run main.go startAsHTTP --port Start as TCP server listening on port 'portNumber'")

}

//validateArgs ensures the cli was given valid input
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		//go exit will exit the application by shutting down the goroutine
		// if you were to use os.exit you might corrupt the data
		runtime.Goexit()
	}
}

func (cli *CommandLine) StartNode(nodeID string) {
	fmt.Printf("Starting Node at %s\n", nodeID)
	network.StartServer(nodeID)
}

func (cli *CommandLine) StartHttp(httpPort string) {
	fmt.Printf("Starting HTTP Server at %s\n", httpPort)
	api.StartServer(httpPort)
}

func (cli *CommandLine) StartgRPC(grpcPort string) {
	fmt.Printf("Starting gRPC Server at %s\n", grpcPort)
	grpcBlockchain.StartServer(grpcPort)
}

//run will start up the command line
func (cli *CommandLine) Run() {
	cli.validateArgs()

	startNodeCmd := flag.NewFlagSet("startAsTCP", flag.ExitOnError)
	startHttpCmd := flag.NewFlagSet("startAsHTTP", flag.ExitOnError)
	startgRPCCmd := flag.NewFlagSet("startAsgRPC", flag.ExitOnError)

	nodeId := startNodeCmd.String("port", "3080", "Start node at this port")
	httpPort := startHttpCmd.String("port", "3080", "Start node at this port")
	blockchainPort := startgRPCCmd.String("blockchainPort", "3080", "Start node at this port")
	grpcPort := startgRPCCmd.String("grpcPort", "9000", "Start node at this port")

	switch os.Args[1] {
	case "startAsTCP":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startAsHTTP":
		err := startHttpCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "startAsgRPC":
		err := startgRPCCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		err1 := startgRPCCmd.Parse(os.Args[3:])
		if err1 != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if startNodeCmd.Parsed() {
		// nodeID := os.Getenv("NODE_ID")
		if *nodeId == "" {
			startNodeCmd.Usage()
			runtime.Goexit()
		}
		port = *nodeId
		cli.StartNode(*nodeId)
	}

	if startHttpCmd.Parsed() {
		if *nodeId == "" {
			startNodeCmd.Usage()
			runtime.Goexit()
		}
		port = *httpPort
		cli.StartHttp(*httpPort)
	}

	if startgRPCCmd.Parsed() {
		if *nodeId == "" {
			startNodeCmd.Usage()
			runtime.Goexit()
		}
		port = *blockchainPort
		cli.StartgRPC(*grpcPort)
	}

}
