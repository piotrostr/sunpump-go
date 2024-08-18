// TODO use the transaction core type and not extention
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/piotrostr/trx/sunpump/processor"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func getClient() *client.GrpcClient {
	// apiKey := os.Getenv("TRONGRID_API_KEY")
	// httpApiURL := "https://api.trongrid.io"
	// grpcURL := "grpc.trongrid.io:50051"

	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
	err := conn.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error starting gRPC client:", err)
		return nil
	}
	return conn
}

func listen(cmd *cobra.Command, args []string) {
	client := getClient()
	if client == nil {
		return
	}

	bch := make(chan *api.BlockExtention)
	stopch := make(chan bool)

	go fetchBlocks(client, bch, stopch)
	go processor.ProcessBlocks(bch, stopch, client)

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	close(stopch)
	fmt.Println("Shutting down...")
}

func getSlot(cmd *cobra.Command, args []string) {
	client := getClient()
	if client == nil {
		return
	}

	block, err := client.GetNowBlock()
	if err != nil {
		fmt.Println("Error getting current slot:", err)
		return
	}

	fmt.Printf("Current slot (block number): %d\n", block.BlockHeader.RawData.Number)
}

func fetchBlocks(client *client.GrpcClient, bch chan<- *api.BlockExtention, stopch <-chan bool) {
	for {
		select {
		case <-stopch:
			return
		default:
			block, err := client.GetNowBlock()
			if err != nil {
				fmt.Println("Error getting block:", err)
				return
			}
			bch <- block
			time.Sleep(2 * time.Second)
		}
	}
}
