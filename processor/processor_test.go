package processor_test

import (
	"testing"

	"github.com/piotrostr/trx/sunpump/processor"
)

func TestDebugTx(t *testing.T) {
	tx, err := processor.ReadTransactionFromFile("mocktx.pb")
	if err != nil {
		t.Fatalf("Failed to read transaction from file: %v", err)
	}

	processor.DebugTx(tx)
}
