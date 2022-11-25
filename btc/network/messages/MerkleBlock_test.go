package messages

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"testing"
)

func TestParseMerkleBlock(t *testing.T) {

	raw, _ := hex.DecodeString("00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635")
	reader := bytes.NewBuffer(raw)
	mb, err := ParseMerkleBlock(reader)

	if err != nil {
		t.Error()
	}

	if mb.Version != 0x20000000 {
		t.Error()
	}

	expectedRoot, _ := hex.DecodeString("ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4")
	utility.ReverseBytes(expectedRoot)
	if !bytes.Equal(expectedRoot, mb.MerkleRoot[:]) {
		t.Error()
	}

	expectedPrevBlock, _ := hex.DecodeString("df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000")
	utility.ReverseBytes(expectedPrevBlock)
	if !bytes.Equal(expectedPrevBlock, mb.PreviousBlock[:]) {
		t.Error()
	}

	expectedTimeStamp := 0x5b837cdc
	if expectedTimeStamp != int(mb.Timestamp) {
		t.Error()
	}

	expectedBits, _ := hex.DecodeString("67d8001a")
	if !bytes.Equal(expectedBits, mb.Bits[:]) {
		t.Error()
	}

	expectedNonce, _ := hex.DecodeString("c157e670")
	if !bytes.Equal(expectedNonce, mb.Nonce[:]) {
		t.Error()
	}

	expectedTotal := 0x00000dbf
	if expectedTotal != int(mb.Total) {
		t.Error()
	}

	expectedHashes := make([][]byte, 10)
	expectedHashes[0], _ = hex.DecodeString("ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a")
	expectedHashes[1], _ = hex.DecodeString("7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d")
	expectedHashes[2], _ = hex.DecodeString("34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2")
	expectedHashes[3], _ = hex.DecodeString("158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cba")
	expectedHashes[4], _ = hex.DecodeString("ee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763ce")
	expectedHashes[5], _ = hex.DecodeString("f8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097")
	expectedHashes[6], _ = hex.DecodeString("c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d")
	expectedHashes[7], _ = hex.DecodeString("6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543")
	expectedHashes[8], _ = hex.DecodeString("d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274c")
	expectedHashes[9], _ = hex.DecodeString("dfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb62261")

	// Reverse them all
	for i := 0; i < 10; i++ {
		utility.ReverseBytes(expectedHashes[i])
	}

	if len(expectedHashes) != len(mb.Hashes) {
		t.Error()
	}

	for i := 0; i < len(mb.Hashes); i++ {
		if !bytes.Equal(expectedHashes[i], mb.Hashes[i][:]) {
			t.Error()
		}
	}

	expectedFlags, _ := hex.DecodeString("b55635")
	if !bytes.Equal(expectedFlags, mb.Flags) {
		t.Error()
	}
}

func TestMerkleBlockIsValid(t *testing.T) {
	raw, _ := hex.DecodeString("00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635")
	reader := bytes.NewBuffer(raw)
	block, err := ParseMerkleBlock(reader)
	if err != nil {
		t.Error()
	}

	if !block.IsValid() {
		t.Error()
	}
}
