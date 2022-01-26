package cli

import (
	"flag"
	"fmt"
	"runtime"

	"github.io/wisysta/nomadcoin/explorer"
	"github.io/wisysta/nomadcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n")
	fmt.Printf("Please use the following commands:\n")
	fmt.Printf("-port=4000: Set the PORT of the server\n")
	fmt.Printf("-mode=rest: Choose between 'html' and 'rest'\n")
	runtime.Goexit()
}

func Start() {
	port := flag.Int("port", 4000, "Set port of the server")

	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)
}
