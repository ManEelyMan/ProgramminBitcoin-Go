package transaction

import (
	"bitcoin-go/collections"
	"bitcoin-go/utility"
	"bytes"
	"testing"
)

func TestOpHash160(t *testing.T) {

	s := collections.NewStack()
	alt := collections.NewStack()

	ctxt := ExecutionContext{Stack: &s, AltStack: &alt}

	s.Push([]byte("hello world"))

	ok := opHash160(&ctxt)
	if !ok {
		t.Error()
	}

	result, _ := s.Pop()
	if !bytes.Equal(utility.HexStringToBigInt("d7d5ee7824ff93f94c3055af9382c86c68b5ca92").Bytes(), result) {
		t.Error()
	}

}
