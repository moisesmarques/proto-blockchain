package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"syscall"

	"github.com/vrecan/death/v3"

	"judge-blockchain/blockchain"
)

var mutex = &sync.RWMutex{}

const (
	protocol        = "tcp"
	version         = 1
	commandLength   = 12
	coinbaseAddress = "1MAD2sfEnSKx2PQP4qPss1sPHgxGEi8guF"
)

var (
	nodeAddress     string
	KnownNodes      = []string{}
	memoryPool      = make(map[string]blockchain.Transaction)
	blocksInTransit = make(map[uint]blockchain.Block)
	chainHeight     = -1
)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}

type GetBlocks struct {
	AddrFrom string
	Height   int
}

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}
type Blocks struct {
	AddrFrom string
	Type     string
	Items    []blockchain.Block
}

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func ExtractCmd(request []byte) []byte {
	return request[:commandLength]
}

func RequestBlocks(chain *blockchain.BlockChain) {
	otherNodes := KnownNodes[1:]
	for _, node := range otherNodes {
		fmt.Printf("Requesting blocks from %s and %d ", node, chainHeight)
		SendGetBlocks(node, chainHeight)
	}
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	SendData(address, request)
}

func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block"), payload...)

	SendData(addr, request)
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	fmt.Printf("sending request to: %s\n", addr)
	// fmt.Printf("sending data: %x \n", data)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}

	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func SendBlocks(address, kind string, items []blockchain.Block) {
	inventory := Blocks{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("blocks"), payload...)

	SendData(address, request)
}

func SendGetBlocks(address string, height int) {
	payload := GobEncode(GetBlocks{nodeAddress, height})
	fmt.Printf("\nGetting blocks from %s and %d \n", address, height)
	request := append(CmdToBytes("getblocks"), payload...)
	// fmt.Printf("\nPayload %x \n", payload)

	SendData(address, request)
}

func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, id})
	request := append(CmdToBytes("getdata"), payload...)

	SendData(address, request)
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	data := Tx{nodeAddress, tnx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(addr, request)
}

func HandleAddr(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)

	}
	fmt.Printf("new request to add known nodes: %s \n", payload)

	// fmt.Printf("previously %d known nodes\n", len(KnownNodes))
	numNewNodes := len(payload.AddrList)
	newNodes := payload.AddrList[:numNewNodes]
	KnownNodes = append(KnownNodes, newNodes...)
	// KnownNodes = append(KnownNodes, payload.AddrF)
	KnownNodes = removeDuplicateNodes(KnownNodes)
	fmt.Printf("there are now %d known nodes\n", len(KnownNodes))
	fmt.Println("********************************************")
	for _, peer := range KnownNodes {
		fmt.Printf("shared node: %s \n", peer)
	}
	fmt.Println("********************************************")
	RequestBlocks(chain)
}
func SendNewBlocksToNetwork(block blockchain.Block) {
	var blocks []blockchain.Block
	blocks = append(blocks, block)

	fmt.Printf("there are now %d known nodes\n", len(KnownNodes))
	fmt.Println("********************************************")
	for _, peer := range KnownNodes {
		SendBlocks(peer, "blocks", blocks)
	}
}

func HandleInv(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocks := payload.Items
		i := 0

		fmt.Println("Adding received blocks!")

		for i <= len(blocks) {
			blockData := blocks[i]
			block := blockchain.Deserialize(blockData)
			mutex.Lock()
			blocksInTransit[uint(block.Height)] = *block
			mutex.Unlock()
			i++
		}
		processBlocksInTransit(chain)
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if memoryPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func HandleBlocks(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Blocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied data with %d %s\n", len(payload.Items), payload.Type)
	fmt.Println("*****************************************************************************")

	if payload.Type == "blocks" {
		blocks := payload.Items
		i := len(blocks)

		fmt.Printf("Received %d brand new blocks!\n", i)
		i--
		for i >= 0 {
			blockData := blocks[i]

			/* if blockData.Height > chainHeight {
				fmt.Println("*******************************")
				fmt.Printf("Adding block height: %d to blocks in transist\n", blockData.Height)
				blocksInTransit[uint(blockData.Height)] = blockData
				fmt.Println("*******************************")
			} */
			mutex.Lock()
			blocksInTransit[uint(blockData.Height)] = blockData
			mutex.Unlock()
			i--
		}
		processBlocksInTransit(chain)
	}
	fmt.Println("*****************************************************************************")

	fmt.Printf("**********************Block Height************************\n\n")
	fmt.Printf("%s chain is at block height: %d \n\n", nodeAddress, chainHeight)
	fmt.Printf("**********************Block Height************************\n\n")

}

func HandleGetBlocks(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%s send a request to fetch blocks from height %d \n", payload.AddrFrom, payload.Height)
	fmt.Println("************************************************")

	// fmt.Println("Request to get blocks: %s")
	if payload.Height < chainHeight {
		fmt.Println("*************Fetching Blocks ***************")
		blocks := chain.GetBlocks(payload.Height)
		if len(blocks) == 0 {
			fmt.Printf("No blocks for address: %s \n", payload.AddrFrom)
		} else {
			fmt.Printf("*************Fetched   %d  blocks, sending to %s\n", len(blocks), payload.AddrFrom)
			SendBlocks(payload.AddrFrom, "blocks", blocks)
		}
	} else {
		fmt.Errorf("Address: %s has a current / more recent version of the chain-No blocks to send\n", payload.AddrFrom)
	}
}

// Read incoming data

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := ioutil.ReadAll(conn)

	defer conn.Close()

	if err != nil {
		log.Panic(err)
	}
	command := BytesToCmd(req[:commandLength])
	fmt.Printf("My adress %s \n", conn.LocalAddr().String())
	fmt.Printf("Received %s command: From %s \n", command, conn.RemoteAddr().String())

	switch command {
	case "addr":
		HandleAddr(req, chain)
	// case "block":
	// 	HandleBlock(req, chain)
	case "inv":
		HandleInv(req, chain)
	case "blocks":
		HandleBlocks(req, chain)
	case "getblocks":
		HandleGetBlocks(req, chain)

	default:
		fmt.Println("Unknown command")
	}

}

func StartServer(nodeID string) {
	blockchain.SetDBPath(nodeID)
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	KnownNodes = append(KnownNodes, nodeAddress)
	// mineAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	if !blockchain.DBexists() {
		fmt.Println("No existing blockchain found, creating one!")
		// these may change depending on intial setup for a peer,
		//Each peer should be setup with atleast the coinbase block
		chain := blockchain.InitBlockChain(coinbaseAddress)
		chain.Database.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Panic(err)
			}
			go HandleConnection(conn, chain)

		}
	} else {
		chain := blockchain.ContinueBlockChain()
		chainHeight = chain.GetBestHeight()
		defer chain.Database.Close()
		go CloseDB(chain)
		fmt.Printf("**********************Block Height************************\n\n")
		fmt.Printf("%s chain is at block height: %d \n\n", nodeAddress, chainHeight)
		fmt.Printf("**********************Block Height************************\n\n")

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Panic(err)
			}
			go HandleConnection(conn, chain)

		}

	}

}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	// fmt.Printf("encoding: %s \n", data)
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Printf("\nencoded %s \n", buff.Bytes())
	return buff.Bytes()
}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}
func removeDuplicateNodes(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			// this will be a check for lenth of of the peer
			if len(entry) > 2 {
				list = append(list, entry)
			}
		}
	}
	return list
}

func processBlocksInTransit(chain *blockchain.BlockChain) {

	for len(blocksInTransit) > 0 {
		currentHeight := uint(chainHeight + 1)
		mutex.RLock()
		value, ok := blocksInTransit[currentHeight]
		mutex.RUnlock()
		if ok {
			fmt.Println("*******************************")
			fmt.Printf("Adding block height: %d to chain\n", value.Height)
			fmt.Printf("Block hash: %x \n", value.Hash)
			fmt.Printf("current height: %d \n", chainHeight)
			if value.Height > chainHeight {
				chain.AddNetworkBlock(&value)
				chainHeight++
			}
			mutex.Lock()
			delete(blocksInTransit, currentHeight)
			mutex.Unlock()
			fmt.Println("*******************************")
		}
	}

}
