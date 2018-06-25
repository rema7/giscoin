package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"giscoin/block"
	"giscoin/concensus"
	"giscoin/transaction"
	"giscoin/wallets"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = ".dbdata/boltdb"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}

func signBlock(block *block.Block) {
	pow := concensus.NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	block.ToSting()
}

func (bc *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := block.NewBlock(transactions, lastHash)
	signBlock(newBlock)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.Db}
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transaction.Transaction {
	var unspentTXs []transaction.Transaction
	spendTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		b := bci.Next()

		for _, tx := range b.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spendTXOs[txID] != nil {
					for _, spentOut := range spendTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spendTXOs[inTxID] = append(spendTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) NewUTXOTransaction(from, to string, amount int) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	_wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := _wallets.GetWallet(from)
	pubKeyHash := wallets.HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("EROR: Not enough funds")
	}

	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := transaction.TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, transaction.TXOutput{amount, pubKeyHash})
	if acc > amount {
		outputs = append(outputs, transaction.TXOutput{acc - amount, pubKeyHash})
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Transaction, error) {
	bci := bc.Iterator()

	for {
		_block := bci.Next()

		for _, tx := range _block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(_block.PrevBlockHash) == 0 {
			break
		}
	}

	return transaction.Transaction{}, nil
}

func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) Verify(tx *transaction.Transaction) bool {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func InitBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		cbtx := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := block.NewGenesisBlock(cbtx)
		signBlock(genesis)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		err = b.Put([]byte("1"), genesis.Hash)
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("Blockchain doesn't exists.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		tip = b.Get([]byte("1"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
