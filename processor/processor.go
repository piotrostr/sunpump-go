// processor processes transactions, blocks, etc, decode/encode etc
package processor

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

func ProcessBlocks(bch <-chan *api.BlockExtention, stopch <-chan bool, client *client.GrpcClient) {
	for {
		select {
		case <-stopch:
			return
		case block := <-bch:
			fmt.Printf("Block %d num txs: %d\n", block.BlockHeader.RawData.Number, len(block.Transactions))
			for _, tx := range block.Transactions {
				contracts := tx.Transaction.RawData.GetContract()
				for _, c := range contracts {
					if c == nil {
						continue
					}
					switch c.Type {
					case core.Transaction_Contract_TriggerSmartContract:
						contract := &core.TriggerSmartContract{}
						err := c.GetParameter().UnmarshalTo(contract)
						if err != nil {
							fmt.Printf("Failed to parse trigger contract: %v\n", err)
							continue
						}
						ca := address.Address(contract.GetContractAddress()).String()
						if ca == "TWekn39KVnSrRML6Sq3xQ5BCJVhYCmoHdb" || ca == "TTfvyrAz86hbZk5iDpKD78pqLGgi8C7AAw" {
							fmt.Printf("Found contract match! %s\n", hex.EncodeToString(tx.Txid))
						}
					default:
						// pass
					}
				}
				// contractAddr := address.Address(c.([]uint8)).String()
			}
		}
	}
}

func WriteTransactionToFile(tx *api.TransactionExtention, filename string) {
	// Serialize the transaction
	data, err := proto.Marshal(tx)
	if err != nil {
		fmt.Printf("Failed to serialize transaction: %v\n", err)
		return
	}

	// Write the serialized data to file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Failed to write transaction to file: %v\n", err)
		return
	}

	fmt.Printf("Transaction written to file: %s\n", filename)
}

func ReadTransactionFromFile(filename string) (*api.TransactionExtention, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Deserialize the transaction
	tx := &api.TransactionExtention{}
	err = proto.Unmarshal(data, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	return tx, nil
}

func DebugTx(tx *api.TransactionExtention) {
	fmt.Printf("Tx: %s\n", tx)
	fmt.Println("Transaction Details:")
	fmt.Printf("Ref Block Bytes: %x\n", tx.Transaction.RawData.RefBlockBytes)
	fmt.Printf("Ref Block Hash: %x\n", tx.Transaction.RawData.RefBlockHash)
	fmt.Printf("Expiration: %v\n", time.Unix(0, tx.Transaction.RawData.Expiration*1000000))

	for _, contract := range tx.Transaction.RawData.Contract {
		fmt.Printf("Contract Type: %v\n", contract.Type)
		switch contract.Type {
		case core.Transaction_Contract_TransferAssetContract:
			var transferContract core.TransferAssetContract
			err := proto.Unmarshal(contract.Parameter.Value, &transferContract)
			if err != nil {
				fmt.Printf("Failed to parse transfer contract: %v\n", err)
				continue
			}
			fmt.Printf("Asset Name: %s\n", transferContract.AssetName)
			fmt.Printf("Owner Address: %x\n", transferContract.OwnerAddress)
			fmt.Printf("To Address: %x\n", transferContract.ToAddress)
			fmt.Printf("Amount: %d\n", transferContract.Amount)
		}
	}

	fmt.Printf("Timestamp: %v\n", time.Unix(0, tx.Transaction.RawData.Timestamp*1000000))
	fmt.Printf("Signature: %x\n", tx.Transaction.Signature[0])

	for _, ret := range tx.Transaction.Ret {
		fmt.Printf("Contract Result: %v\n", ret.ContractRet)
	}

}

func AssetName(tx *api.TransactionExtention) string {
	fmt.Printf("Tx contracts len: %d\n", len(tx.Transaction.RawData.Contract))
	for _, contract := range tx.Transaction.RawData.Contract {
		if contract.Type == core.Transaction_Contract_TransferAssetContract {
			var transferContract core.TransferAssetContract
			err := proto.Unmarshal(contract.Parameter.Value, &transferContract)
			if err != nil {
				fmt.Printf("Failed to parse transfer contract: %v\n", err)
				continue
			}
			return string(transferContract.AssetName)
		}
	}
	return ""
}
