// TODO use the transaction core type and not extention
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/fatih/structs"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/piotrostr/trx/sunpump/processor"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
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

var getTxCmd = &cobra.Command{
	Use:   "getTx",
	Short: "Get a transaction by ID",
	Run:   getTx,
}

func init() {
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(getSlotCmd)
	rootCmd.AddCommand(getTxCmd)
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

func getTx(cmd *cobra.Command, args []string) {
	client := getClient()
	if client == nil {
		return
	}

	tx, err := client.GetTransactionByID(args[0])
	if err != nil {
		fmt.Println("Error getting transaction:", err)
		return
	}

	err = parseContract(tx.RawData)
	if err != nil {
		fmt.Println("Error parsing contract:", err)
		return
	}
}

func parseContract(txRaw *core.TransactionRaw) error {
	for _, contract := range txRaw.GetContract() {
		var c interface{}
		switch contract.Type {
		case core.Transaction_Contract_AccountCreateContract:
			c = &core.AccountCreateContract{}
		case core.Transaction_Contract_TransferContract:
			c = &core.TransferContract{}
		case core.Transaction_Contract_TransferAssetContract:
			c = &core.TransferAssetContract{}
		case core.Transaction_Contract_VoteWitnessContract:
			c = &core.VoteWitnessContract{}
		case core.Transaction_Contract_WitnessCreateContract:
			c = &core.WitnessCreateContract{}
		case core.Transaction_Contract_WitnessUpdateContract:
			c = &core.WitnessUpdateContract{}
		case core.Transaction_Contract_AssetIssueContract:
			c = &core.AssetIssueContract{}
		case core.Transaction_Contract_ParticipateAssetIssueContract:
			c = &core.ParticipateAssetIssueContract{}
		case core.Transaction_Contract_AccountUpdateContract:
			c = &core.AccountUpdateContract{}
		case core.Transaction_Contract_FreezeBalanceContract:
			c = &core.FreezeBalanceContract{}
		case core.Transaction_Contract_UnfreezeBalanceContract:
			c = &core.UnfreezeBalanceContract{}
		case core.Transaction_Contract_WithdrawBalanceContract:
			c = &core.WithdrawBalanceContract{}
		case core.Transaction_Contract_UnfreezeAssetContract:
			c = &core.UnfreezeAssetContract{}
		case core.Transaction_Contract_UpdateAssetContract:
			c = &core.UpdateAssetContract{}
		case core.Transaction_Contract_ProposalCreateContract:
			c = &core.ProposalCreateContract{}
		case core.Transaction_Contract_ProposalApproveContract:
			c = &core.ProposalApproveContract{}
		case core.Transaction_Contract_ProposalDeleteContract:
			c = &core.ProposalDeleteContract{}
		case core.Transaction_Contract_SetAccountIdContract:
			c = &core.SetAccountIdContract{}
		case core.Transaction_Contract_CustomContract:
			return fmt.Errorf("proto unmarshal any: %s", "customContract")
		case core.Transaction_Contract_CreateSmartContract:
			c = &core.CreateSmartContract{}
		case core.Transaction_Contract_TriggerSmartContract:
			c = &core.TriggerSmartContract{}
		case core.Transaction_Contract_GetContract:
			return fmt.Errorf("proto unmarshal any: %s", "getContract")
		case core.Transaction_Contract_UpdateSettingContract:
			c = &core.UpdateSettingContract{}
		case core.Transaction_Contract_ExchangeCreateContract:
			c = &core.ExchangeCreateContract{}
		case core.Transaction_Contract_ExchangeInjectContract:
			c = &core.ExchangeInjectContract{}
		case core.Transaction_Contract_ExchangeWithdrawContract:
			c = &core.ExchangeWithdrawContract{}
		case core.Transaction_Contract_ExchangeTransactionContract:
			c = &core.ExchangeTransactionContract{}
		case core.Transaction_Contract_UpdateEnergyLimitContract:
			c = &core.UpdateEnergyLimitContract{}
		case core.Transaction_Contract_AccountPermissionUpdateContract:
			c = &core.AccountPermissionUpdateContract{}
		case core.Transaction_Contract_ClearABIContract:
			c = &core.ClearABIContract{}
		case core.Transaction_Contract_UpdateBrokerageContract:
			c = &core.UpdateBrokerageContract{}
		case core.Transaction_Contract_ShieldedTransferContract:
			c = &core.ShieldedTransferContract{}
		case core.Transaction_Contract_MarketSellAssetContract:
			c = &core.MarketSellAssetContract{}
		case core.Transaction_Contract_MarketCancelOrderContract:
			c = &core.MarketCancelOrderContract{}
		case core.Transaction_Contract_FreezeBalanceV2Contract:
			c = &core.FreezeBalanceV2Contract{}
		case core.Transaction_Contract_UnfreezeBalanceV2Contract:
			c = &core.UnfreezeBalanceV2Contract{}
		case core.Transaction_Contract_WithdrawExpireUnfreezeContract:
			c = &core.WithdrawExpireUnfreezeContract{}
		case core.Transaction_Contract_DelegateResourceContract:
			c = &core.DelegateResourceContract{}
		case core.Transaction_Contract_UnDelegateResourceContract:
			c = &core.UnDelegateResourceContract{}
		default:
			return fmt.Errorf("proto unmarshal any")
		}

		if err := contract.GetParameter().UnmarshalTo(c.(protoreflect.ProtoMessage)); err != nil {
			return fmt.Errorf("proto unmarshal any: %+w", err)
		}

		// hrc := parseContractHumanReadable(structs.Map(c))
		json, err := json.Marshal(parseContractHumanReadable(structs.Map(c)))
		if err != nil {
			return fmt.Errorf("json marshal: %+v", err)
		}
		fmt.Println(string(json))
	}
	return nil
}

func parseContractHumanReadable(ck map[string]interface{}) map[string]interface{} {
	// Addresses fields
	addresses := map[string]bool{
		"OwnerAddress":    true,
		"ReceiverAddress": true,
		"ToAddress":       true,
		"ContractAddress": true,
	}
	for f, d := range ck {
		if strings.HasPrefix(f, "XXX_") {
			delete(ck, f)
		}

		// convert addresses
		if addresses[f] {
			ck[f] = address.Address(d.([]uint8)).String()
		}
	}

	if v, ok := ck["Votes"]; ok {
		votes := make(map[string]int64)
		for _, d := range v.([]interface{}) {
			dP := d.(map[string]interface{})
			votes[address.Address(dP["VoteAddress"].([]uint8)).String()] = dP["VoteCount"].(int64)
		}
		ck["Votes"] = votes
	}

	return ck
}

func fetchBlocks(client *client.GrpcClient, bch chan<- *api.BlockExtention, stopch <-chan bool) {
	for {
		block, err := client.GetNowBlock()
		if err != nil {
			fmt.Println("Error getting block:", err)
			return
		}
		bch <- block
		if <-stopch {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
