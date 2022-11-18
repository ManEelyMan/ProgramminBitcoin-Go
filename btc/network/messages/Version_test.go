package messages

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestVersionMessage(t *testing.T) {

	v := NewDefaultVersionMessage()

	writer := bytes.NewBuffer(make([]byte, 0))
	v.Serialize(writer)

	expected, _ := hex.DecodeString("7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000")

	actual := writer.Bytes()
	if !bytes.Equal(expected, actual) {
		t.Error()
	}
}
