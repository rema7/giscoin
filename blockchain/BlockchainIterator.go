package blockchain

import (
	"giscoin/block"
	"giscoin/utils"
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bi *BlockchainIterator) Next() *block.Block {
	var _block *block.Block

	err := bi.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bi.currentHash)
		_block, _ = utils.DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bi.currentHash = _block.PrevBlockHash
	return _block
}
