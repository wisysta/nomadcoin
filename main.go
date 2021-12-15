package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevhash string
}

func main() {
	genesisBlock := block{"Genesis Block", "", ""}
	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevhash))
	hexHash := fmt.Sprintf("%x", hash)
	genesisBlock.hash = hexHash

	fmt.Println(genesisBlock)
}
