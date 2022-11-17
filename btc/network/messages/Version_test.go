package messages

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestVersionMessage(t *testing.T) {

	v := Version{
		Version:                70015,
		SenderNetworkPort:      [2]byte{0x20, 0x8D},
		ReceiverNetworkPort:    [2]byte{0x20, 0x8D},
		Timestamp:              0,
		Nonce:                  [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		UserAgent:              "/programmingbitcoin:0.1/",
		ReceiverNetworkAddress: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0},
		SenderNetworkAddress:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0},
	}

	writer := bytes.NewBuffer(make([]byte, 0))
	v.Serialize(writer)

	expected, _ := hex.DecodeString("7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000")

	actual := writer.Bytes()
	if !bytes.Equal(expected, actual) {
		t.Error()
	}
}
