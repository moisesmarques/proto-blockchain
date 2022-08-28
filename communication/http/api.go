package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"judge-blockchain/blockchain"
	network "judge-blockchain/communication/tcp"
	"judge-blockchain/wallet"

	"github.com/gorilla/mux"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the PrimeChain API!\n")
	fmt.Fprintf(w, "Endpoints\n")
	fmt.Fprintf(w, "*************************************************\n")
	fmt.Fprintf(w, "Usage: \n")

	fmt.Fprintf(w, "POST /chain body: {address:''} creates a blockchain and rewards the mining fee\n")
	fmt.Fprintf(w, "GET /chain - Prints the blocks in the chain\n")

	fmt.Fprintf(w, "POST /wallets - Creates a new wallet\n")
	fmt.Fprintf(w, "GET /wallets - Lists the addresses in the wallet file with their balances\n")
	fmt.Fprintf(w, "POST /sendTokens body: {from:FROM,to:TO,amount:AMOUNT} - Send amount of coins from one address to another\n")
	fmt.Fprintf(w, "POST /newActionResult body: {our validator results schema} send validator action results and save in blockchain \n")

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func createChain(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// fmt.Println(vars["wallet"])
	address := r.FormValue("wallet")
	if !wallet.ValidateAddress(address) {
		w.WriteHeader(http.StatusOK)
		// fmt.Fprintf(w, "Wallet is invalid")
	}
	if blockchain.DBexists() {
		fmt.Fprintf(w, "Blockchain exists\n")
	}
	fmt.Fprintf(w, "*****************************Printing chain data********************************\n")
	newChain := blockchain.InitBlockChain(address)
	newChain.Database.Close()
	// fmt.Println("Finished creating chain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "*************************************************************\n")
	printChain(w, r)
}
func sendTokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("wallet"))
	from := r.FormValue("from")
	to := r.FormValue("to")
	// valid := wallet.ValidateAddress(from) && wallet.ValidateAddress(to)

	if amount, err := strconv.Atoi(r.FormValue("amount")); err == nil && len(from) != 0 && len(to) != 0 {

		chain := blockchain.ContinueBlockChain()
		defer chain.Database.Close()

		wallets, err := wallet.CreateWallets()
		if err != nil {
			log.Panic(err)
		}
		wallet := wallets.GetWallet(from)

		tx := blockchain.NewTransaction(&wallet, to, amount, chain)

		chain.AddBlock([]*blockchain.Transaction{tx})
		// send it to the network
		// network.SendNewBlocksToNetwork(*block)
		SendNewBlockToNetwork(*chain)
		fmt.Println("Added Transaction Success!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Tokens send successfully \n")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Inputs to, form and amount must be provided \n")
	}
}

func SendNewBlockToNetwork(chain blockchain.BlockChain) {
	// var lastHash []byte

	// err := chain.Database.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get([]byte("lh"))
	// 	blockchain.Handle(err)

	// 	err = item.Value(func(val []byte) error {
	// 		lastHash = val

	// 		return nil
	// 	})
	// 	blockchain.Handle(err)

	// 	return err
	// })
	// blockchain.Handle(err)

	// iterator := chain.Iterator()
	// lastBlock := iterator.Next()

	// fmt.Println("**********************")
	// fmt.Printf("Hash %s", lastBlock.Hash)
	// // fmt.Printf("Last Hash %s", chain.LastHash)
	// lastBlock := blockchain.Deserialize(lastHash)

	iterator := chain.Iterator()
	lastBlock := iterator.Next()

	var blocks []blockchain.Block
	blocks = append(blocks, *lastBlock)

	fmt.Printf("there are now %d known nodes\n", len(network.KnownNodes))
	fmt.Println("********************************************")
	for _, peer := range network.KnownNodes {
		network.SendBlocks(peer, "blocks", blocks)
		// network.SendBlock(peer, lastBlock)
	}
}

func deleteChain(w http.ResponseWriter, r *http.Request) {
	if !blockchain.DBexists() {
		w.WriteHeader(http.StatusOK)
		fmt.Println("No blockchain found")
		fmt.Fprintf(w, "No blockchain found")
		runtime.Goexit()
	} else {
		blockchain.DeleteBlockChain()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Chain deleted")
	}
}
func listWallets(w http.ResponseWriter, r *http.Request) {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()
	w.WriteHeader(http.StatusOK)
	for _, address := range addresses {
		fmt.Println(address)
		fmt.Fprintf(w, "Address: %s  Balance: %d\n", address, getBalance(address))
	}
}
func getBalance(address string) int {
	fmt.Println(address)
	// if !wallet.ValidateAddress(address) {
	// 	fmt.Printf("Address is not Valid")
	// 	return 0
	// }
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := chain.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
	return balance

}
func createCWallet(w http.ResponseWriter, r *http.Request) {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "New address is: %s\n", address)

}
func printChain(w http.ResponseWriter, r *http.Request) {
	if !blockchain.DBexists() {
		w.WriteHeader(http.StatusOK)
		fmt.Println("No blockchain found, please create one first")
		fmt.Fprintf(w, "No blockchain found, please create one first")
		// runtime.Goexit()
	} else {

		chain := blockchain.ContinueBlockChain()
		defer chain.Database.Close()
		iterator := chain.Iterator()

		w.WriteHeader(http.StatusOK)
		// fmt.Fprintf(w, "API is up and running")

		for {
			block := iterator.Next()
			// fmt.Printf("Previous hash: %x\n", block.PrevHash)
			fmt.Fprintf(w, "Previous hash: %x\n", block.PrevHash)
			fmt.Fprintf(w, "hash: %x\n", block.Hash)
			// fmt.Printf("hash: %x\n", block.Hash)
			pow := blockchain.NewProofOfWork(block)
			fmt.Fprintf(w, "Pow: %s\n", strconv.FormatBool(pow.Validate()))
			for _, tx := range block.Transactions {
				fmt.Println(tx)
				fmt.Fprintf(w, " %s\n", tx)
			}
			fmt.Fprintf(w, "Height: %d\n", block.Height)
			fmt.Fprintf(w, "Nonce: %d\n", block.Nonce)
			fmt.Fprintf(w, "Timestamp: %d\n", block.Timestamp)
			fmt.Println("*************************************************")
			fmt.Fprintf(w, "*******************************************************************************\n")
			// This works because the Genesis block has no PrevHash to point to.
			if len(block.PrevHash) == 0 {
				break
			}
		}
	}
}

func newActionResults(w http.ResponseWriter, r *http.Request) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	actionResul := blockchain.ValidatorActionResult{} //initialize empty user
	err := json.NewDecoder(r.Body).Decode(&actionResul)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "err: %x\n", err.Error())
	}
	fmt.Println(actionResul)

	tx := blockchain.NewActionResult(chain, actionResul)

	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")

	iterator := chain.Iterator()

	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "API is up and running")

	/* for {
		block := iterator.Next()
		fmt.Fprintf(w, "Previous hash: %x\n", block.PrevHash)
		fmt.Fprintf(w, "hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Fprintf(w, "Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevHash) == 0 {
			break
		}
	}
	*/

	// broadcast the block to other nodes
	lastBlock := iterator.Next()

	fmt.Fprintf(w, "Previous hash: %x\n", lastBlock.PrevHash)
	fmt.Fprintf(w, "hash: %x\n", lastBlock.Hash)
	/* pow := blockchain.NewProofOfWork(lastBlock)
	fmt.Fprintf(w, "Pow: %s\n", strconv.FormatBool(pow.Validate())) */

	var blocks []blockchain.Block
	blocks = append(blocks, *lastBlock)
	for _, peer := range network.KnownNodes {
		network.SendBlocks(peer, "blocks", blocks)
		//network.SendBlock(peer, lastBlock)
	}
}

func listPeers(w http.ResponseWriter, r *http.Request) {
	peers := network.KnownNodes
	for _, peer := range peers {
		fmt.Println(peer)
		fmt.Fprintf(w, " %s\n", peer)
	}
	fmt.Println("*************************************************")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All addresses printed")
}

func addPeer(w http.ResponseWriter, r *http.Request) {
	host := r.FormValue("host")
	port := r.FormValue("port")

	nodeExist := false

	// peers := network.KnownNodes
	for _, peer := range network.KnownNodes {
		if peer == fmt.Sprintf("%s:%s", host, port) {
			// fmt.Println(peer)
			fmt.Fprintf(w, " %s Exists\n", peer)
			nodeExist = true
		}
	}

	if !nodeExist {
		newAdrss := fmt.Sprintf("%s:%s", host, port)
		network.KnownNodes = append(network.KnownNodes, newAdrss)
		fmt.Printf("known addresses to %s \n", network.KnownNodes)
	}

	for _, peer := range network.KnownNodes {
		network.SendAddr(peer)
		fmt.Printf("send addresses to %s \n", peer)
	}
	// broadcast new address to peers

	// fmt.Println("send tx")
	fmt.Println("*************************************************")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Peer added\n")
}

func StartServer(port string) {

	blockchain.SetDBPath(port)

	thisNode := fmt.Sprintf("localhost:%s", port)
	network.KnownNodes = append(network.KnownNodes, thisNode)
	fmt.Printf("Known nodes: %s", network.KnownNodes)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homepage).Methods("GET")
	router.HandleFunc("/chain", printChain).Methods("GET")
	router.HandleFunc("/chain", createChain).Methods("POST")
	router.HandleFunc("/chain", deleteChain).Methods("DELETE")
	router.HandleFunc("/wallets", listWallets).Methods("Get")
	router.HandleFunc("/wallets", createCWallet).Methods("POST")

	// router.HandleFunc("/transaction", createTransaction).Methods("POST")
	// router.HandleFunc("/transactions", getTransactions).Methods("GET")
	// router.HandleFunc("/balance", walletBalance).Methods("Get")
	router.HandleFunc("/sendTokens", sendTokens).Methods("POST")
	router.HandleFunc("/sendTokens", sendTokens).Methods("POST")
	router.HandleFunc("/newActionResult", newActionResults).Methods("POST")

	router.HandleFunc("/peers", listPeers).Methods("GET")
	router.HandleFunc("/peers", addPeer).Methods("POST")
	router.HandleFunc("/health", HealthCheck).Methods("GET")

	http.ListenAndServe(":"+port, router)
}
