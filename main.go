// TODO use the transaction core type and not extention
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/distributed-lab/tron-sdk/tron_api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/piotrostr/trx/sunpump/processor"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sunpump",
	Short: "Sunpump is a tool for processing TRON blocks",
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start listening for and processing TRON blocks",
	Run:   listen,
}

var getSlotCmd = &cobra.Command{
	Use:   "getSlot",
	Short: "Get the current slot number",
	Run:   getSlot,
}

func init() {
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(getSlotCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listen(cmd *cobra.Command, args []string) {
	apiKey := os.Getenv("TRONGRID_API_KEY")
	httpApiURL := "https://api.trongrid.io"
	grpcURL := "grpc.trongrid.io:50051"

	client := tron_api.NewTronClient(httpApiURL, grpcURL, apiKey)

	bnch := make(chan int64)
	bch := make(chan *api.BlockExtention)
	stopch := make(chan bool)

	go produceBlockNumbers(client, bnch, stopch)
	go fetchBlocks(client, bnch, bch, stopch)
	go processor.ProcessBlocks(bch, stopch)

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	close(stopch)
	fmt.Println("Shutting down...")
}

func getSlot(cmd *cobra.Command, args []string) {
	apiKey := os.Getenv("TRONGRID_API_KEY")
	httpApiURL := "https://api.trongrid.io"
	grpcURL := "grpc.trongrid.io:50051"

	client := tron_api.NewTronClient(httpApiURL, grpcURL, apiKey)

	blockNumber, err := client.GetNowBlock()
	if err != nil {
		fmt.Println("Error getting current slot:", err)
		return
	}

	fmt.Printf("Current slot (block number): %d\n", blockNumber)
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
