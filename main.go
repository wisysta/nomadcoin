package main

import (
	"github.io/wisysta/nomadcoin/cli"
	"github.io/wisysta/nomadcoin/db"
)

func main() {
	// blockchain.Blockchain().AddBlock("First")
	// blockchain.Blockchain().AddBlock("Second")
	// blockchain.Blockchain()
	defer db.Close()
	cli.Start()
}
