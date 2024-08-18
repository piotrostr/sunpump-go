// processor processes transactions, blocks, etc, decode/encode etc
package processor

import (
	"fmt"
	"os"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

func ProcessBlocks(bch <-chan *api.BlockExtention, stopch <-chan bool) {
	for {
		select {
		case <-stopch:
			return
		case block := <-bch:
			fmt.Printf("Block %d num txs: %d\n", block.BlockHeader.RawData.Number, len(block.Transactions))
			for i, tx := range block.Transactions {
				// Write the transaction to file
				filename := fmt.Sprintf("transaction_%d_%d.pb", block.BlockHeader.RawData.Number, i)
				WriteTransactionToFile(tx, filename)
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

func DecodeTx() {}
