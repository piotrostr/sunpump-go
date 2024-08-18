// TODO use the transaction core type and not extention
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/distributed-lab/tron-sdk/tron_api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

func main() {
	apiKey := os.Getenv("TRONGRID_API_KEY")
	httpApiURL := "https://api.trongrid.io"
	grpcURL := "grpc.trongrid.io:50051"

	client := tron_api.NewTronClient(httpApiURL, grpcURL, apiKey)

	bnch := make(chan int64)
	bch := make(chan *api.BlockExtention)
	stopch := make(chan bool)

	go produceBlockNumbers(client, bnch, stopch)
	go fetchBlocks(client, bnch, bch, stopch)
	go processBlocks(bch, stopch)

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	close(stopch)
	fmt.Println("Shutting down...")
}

func produceBlockNumbers(client *tron_api.TronClient, bnch chan<- int64, stopch <-chan bool) {
	for {
		select {
		case <-stopch:
			return
		default:
			blockNumber, err := client.GetNowBlock()
			if err != nil {
				fmt.Println("Error getting latest block number:", err)
				continue
			}
			bnch <- blockNumber
		}
	}
}

func fetchBlocks(client *tron_api.TronClient, bnch <-chan int64, bch chan<- *api.BlockExtention, stopch <-chan bool) {
	for {
		select {
		case <-stopch:
			return
		case blockNumber := <-bnch:
			block, err := client.GetBlockByNum(blockNumber)
			if err != nil {
				fmt.Println("Error fetching block:", err)
				continue
			}
			bch <- block
		}
	}
}

func processBlocks(bch <-chan *api.BlockExtention, stopch <-chan bool) {
	for {
		select {
		case <-stopch:
			return
		case block := <-bch:
			fmt.Printf("Block %d num txs: %d", block.BlockHeader.RawData.Number, len(block.Transactions))
			for _, tx := range block.Transactions {
				fmt.Printf("Transaction: %+v\n", tx)
			}
		}
	}
}
