package messages

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestGetHeadersMessageSerialize(t *testing.T) {

	var blockHex [32]byte
	arr, _ := hex.DecodeString("0000000000000000001237f46acddf58578a37e213d2a6edc4884a2fcad05ba3")
	copy(blockHex[:], arr)

	gh := GetHeaders{
		Version:    70015,
		NumHashes:  1,
		StartBlock: blockHex,
	}

	writer := bytes.NewBuffer(make([]byte, 0))
	gh.Serialize(writer)
	actual := writer.Bytes()

	expected, _ := hex.DecodeString("7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	if !bytes.Equal(expected, actual) {
		t.Error()
	}
}
