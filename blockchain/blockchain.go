package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	utils "judge-blockchain/utils"

	"github.com/dgraph-io/badger"
)

var (
	dbPath = "./database/blockchain/nodes/3080"

	// This can be used to verify that the blockchain exists
	dbFile = "./database/blockchain/nodes/3080/MANIFEST"

	// This is arbitrary data for our genesis block
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func SetDBPath(port string) {
	dbPath = "./database/blockchain/nodes/" + port
	// This can be used to verify that the blockchain exists
	dbFile = "./database/blockchain/nodes/" + port + "/MANIFEST"

	if !DBexists() {
		utils.CopyDir("./database/blockchain/basechain", dbPath)
	}
}

func ContinueBlockChain() *BlockChain {

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	// db, err := badger.Open(opts)
	db, err := openDB(dbPath, opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

func DeleteBlockChain() {

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)

	err = db.DropAll()
	db.Close()
	Handle(err)
	// runtime.Goexit()
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("blockchain already exists")
		blockchain := ContinueBlockChain()
		return blockchain
	}

	opts := badger.DefaultOptions(dbPath)
	// db, err := badger.Open(opts)
	db, err := openDB(dbPath, opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err

	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val

			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	iterator := chain.Iterator()
	lastBlock := iterator.Next()

	fmt.Println("**********************")
	fmt.Printf("Hash %s", lastBlock.Hash)
	// fmt.Printf("Last Hash %s", chain.LastHash)
	// lastBlock := Deserialize(lastHash)
	fmt.Printf("Last Block Height %d", lastBlock.Height)
	newBlock := CreateBlock(transactions, lastHash, lastBlock.Height+1)

	err = chain.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
	return newBlock

	// var lastHash []byte

	// for _, tx := range transactions {
	// 	if chain.VerifyTransaction(tx) != true {
	// 		log.Panic("Invalid Transaction")
	// 	}
	// }

	// err := chain.Database.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get([]byte("lh"))
	// 	// lastBlockData, _ := item.Value()
	// 	fmt.Printf(err.Error())
	// 	Handle(err)
	// 	// lastHash, err = item.Value()
	// 	err = item.Value(func(val []byte) error {
	// 		lastHash = val
	// 		return nil
	// 	})

	// 	return err
	// })
	// Handle(err)

	// // TO DO add height

	// // lastBlockData, _ := item.Value()

	// lastBlock := Deserialize(lastHash)

	// // if block.Height > lastBlock.Height {
	// // 	err = txn.Set([]byte("lh"), block.Hash)
	// // 	Handle(err)
	// // 	chain.LastHash = block.Hash
	// // }
	// newBlock := CreateBlock(transactions, lastHash, lastBlock.Height+1)

	// err = chain.Database.Update(func(txn *badger.Txn) error {
	// 	err := txn.Set(newBlock.Hash, newBlock.Serialize())
	// 	Handle(err)
	// 	err = txn.Set([]byte("lh"), newBlock.Hash)

	// 	chain.LastHash = newBlock.Hash

	// 	return err
	// })
	// Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}

	return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
func (chain *BlockChain) GetBestHeight() int {
	// var lastBlock Block
	// var lastHash []byte

	iterator := chain.Iterator()
	lastBlock := iterator.Next()

	// err := chain.Database.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get([]byte("lh"))
	// 	Handle(err)

	// 	err = item.Value(func(val []byte) error {
	// 		lastHash = val
	// 		return nil
	// 	})

	// 	return err
	// })
	// // 	err = item.Value(func(val []byte) error {
	// // 		lastHash = val
	// // 		return nil
	// // 	})

	// // 	lastBlock = *Deserialize(lastHash)

	// // 	return nil
	// // })
	// Handle(err)

	// lastBlock = *Deserialize(lastHash)

	// // err = item.Value(func(val []byte) error {
	// // 	lastHash = val
	// // 	return nil
	// // })

	return lastBlock.Height
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	lastBlock := Deserialize(lastHash)

	lastHeight = lastBlock.Height

	Handle(err)

	// newBlock := CreateBlock(transactions, lastHash)
	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
	return newBlock
}

func (chain *BlockChain) GetBlockHashes(height int) [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()
	thisChainHeight := chain.GetBestHeight()
	counter := thisChainHeight - height

	for counter > 0 {
		block := iter.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
		counter--
	}

	return blocks
}

func (chain *BlockChain) GetBlocks(height int) []Block {
	var blocks []Block

	iter := chain.Iterator()
	thisChainHeight := chain.GetBestHeight()
	counter := thisChainHeight - height

	for counter > 0 {
		block := iter.Next()

		blocks = append(blocks, *block)

		if len(block.PrevHash) == 0 {
			break
		}
		counter--
	}

	return blocks
}
func (chain *BlockChain) AddNetworkBlock(block *Block) {
	if chain.isValidBlock(block) {
		err := chain.Database.Update(func(txn *badger.Txn) error {
			if _, err := txn.Get(block.Hash); err == nil {
				return nil
			}

			blockData := block.Serialize()
			err := txn.Set(block.Hash, blockData)
			Handle(err)

			err = txn.Set([]byte("lh"), block.Hash)
			Handle(err)
			chain.LastHash = block.Hash

			return nil
		})
		Handle(err)
	} else {
		fmt.Printf("Block is invalid %d \n", block.Height)
		// fmt.Errorf("Block is invalid %d \n", block.Height)
	}
}
func (chain *BlockChain) isValidBlock(block *Block) bool {

	iterator := chain.Iterator()
	lastBlock := iterator.Next()

	//pow := NewProofOfWork(block)
	//proof := pow.Validate()
	hashCheck := reflect.DeepEqual(lastBlock.Hash, block.PrevHash)
	heightCheck := lastBlock.Height < block.Height
	timestampCheck := lastBlock.Timestamp < block.Timestamp
	fmt.Println("*******************************")
	// fmt.Printf("Last Chain Hash: %x \n and this block lastHash: %x \n", lastBlock.Hash, block.PrevHash)
	fmt.Printf("Last Block Height check: %d is <  %d : %t \n", lastBlock.Height, block.Height, heightCheck)
	fmt.Printf("Timestamps check: %d is < %d : %t \n", lastBlock.Timestamp, block.Timestamp, timestampCheck)
	fmt.Printf("Previous Hash check: %x \nCompare lastHash: %x  \nCheck is: %t \n", lastBlock.Hash, block.PrevHash, hashCheck)
	//fmt.Printf("Proof of work: %t \n\n", proof)
	fmt.Println("*******************************")
	// return proof && heightCheck && hashCheck && timestampCheck
	/* the proof of work check seems not to work when validating a network block */
	return heightCheck && hashCheck && timestampCheck
}
