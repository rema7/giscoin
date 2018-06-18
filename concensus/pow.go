package concensus

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"giscoin/block"
	"giscoin/utils"
	"math"
	"math/big"
)

const targetBits = 20

type ProofOfWork struct {
	block  *block.Block
	target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		utils.IntToHex(pow.block.Timestamp),
		utils.IntToHex(int64(targetBits)),
		utils.IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
		fmt.Print("\n\n")
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hasInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hasInt.SetBytes(hash[:])

	return hasInt.Cmp(pow.target) == -1
}

func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}
