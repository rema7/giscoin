package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"giscoin/block"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeBlock(d []byte) (*block.Block, error) {
	var b block.Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)

	return &b, err
}
