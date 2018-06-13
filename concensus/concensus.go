package concensus

import "giscoin/block"

func SingBlock(newBlock *block.Block) (int, []byte) {
	pow := NewProofOfWork(newBlock)
	return pow.Run()
}
