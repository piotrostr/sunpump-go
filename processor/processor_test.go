package processor_test

// // func TestProcessAndDecodeTransaction(t *testing.T) {
// // 	// Assuming you have a way to generate test blocks and transactions
// // 	bch := make(chan *api.BlockExtention)
// // 	stopch := make(chan bool)
// //
// // 	go ProcessBlocks(bch, stopch)
// //
// // 	// Send a test block
// // 	testBlock := generateTestBlock() // You need to implement this function
// // 	bch <- testBlock
// //
// // 	// Wait for processing
// // 	time.Sleep(time.Second)
// //
// // 	// Stop processing
// // 	stopch <- true
// //
// // 	// Now, test decoding
// // 	filename := fmt.Sprintf("transaction_%d_0.pb", testBlock.BlockHeader.RawData.Number)
// // 	DecodeTx(filename)
// //
// // 	// Add your assertions here
// // 	// ...
// // }
