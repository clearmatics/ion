// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli_test

// func Test_EncodeBlock(t *testing.T) {
// 	// read a fake block
// 	raw, _ := ioutil.ReadFile("../block.json")
//
// 	const expectedRlpHex = "f90256a0ad34f0f919e4b06b18b0c674b8b9f6738a4878c76e837c8f31a2079f21dced1ca01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0db37435caa1fca7e1aa5b4da1c69fdf1d127232519eb3b1b5069825e6c62f5dca056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a83fd82da80845b28beecb861d78301080a846765746887676f312e392e33856c696e75780000000000000000e0ac79c5577889dfb5745ace9c5dfebe1a11bb19ced9b98b427e7bd4c85765ce17154e658440915743ec442fb64756483bc592616754d13a3c62fce5a56ac9f501a00000000000000000000000000000000000000000000000000000000000000000880000000000000000"
//
// 	var marshalledBlock cli.Header
// 	json.Unmarshal(raw, &marshalledBlock)
//
// 	// Now RLP encode the block
// 	blockInterface := cli.GenerateInterface(marshalledBlock)
// 	fmt.Printf("%+v\n", marshalledBlock.Extra)
// 	hash := cli.EncodeBlock(blockInterface)
// 	assert.Equal(t, expectedRlpHex, hex.EncodeToString(hash))
// }