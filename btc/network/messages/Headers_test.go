package messages

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHeadersMessageParse(t *testing.T) {

	msg, _ := hex.DecodeString("0200000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670000000002030eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000768b89f07044e6130ead292a3f51951adbd2202df447d98789339937fd006bd44880835b67d8001ade09204600")
	reader := bytes.NewBuffer(msg)

	tmp, err := parseHeadersMessage(reader)
	if err != nil {
		t.Error()
	}

	hdr := tmp.(Headers)
	if len(hdr.Blocks) != 2 {
		t.Error()
	}

	// Check for parsing errors.
	if hdr.Blocks[0].Version != 536870912 {
		t.Error()
	}

	if hdr.Blocks[1].Version != 536870912 {
		t.Error()
	}
}
