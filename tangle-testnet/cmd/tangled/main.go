package main

import (
	"flag"
	"fmt"
	"github.com/ICFTProtocol/tangle-testnet/internal/api"
	"github.com/ICFTProtocol/tangle-testnet/internal/store"
	"log"
)

func main() {
	data := flag.String("data", "./data/tangle.json", "state file")
	httpAddr := flag.String("http", ":8080", "HTTP listen address")
	chainID := flag.String("chain-id", "tangle-testnet-1", "chain id")
	flag.Parse()
	s, err := store.New(*data, *chainID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TANGLE testnet alpha — chain=%s height=%d http=%s\n", s.State.ChainID, s.State.Height, *httpAddr)
	api.Start(s, *httpAddr)
}
