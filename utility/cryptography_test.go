package utility

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBase58(t *testing.T) {

	addr := "mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf"
	h160, ok := DecodeBase58(addr)
	if !ok {
		t.Error()
	}

	expected, err := hex.DecodeString("507b27411ccf7f16f10297de6cef3f291623eddf")
	if err != nil {
		t.Error()
	}

	if !bytes.Equal(h160, expected) {
		t.Error()
	}

	// Round-trip it.
	res := H160ToP2PKHAddress(h160, true)
	if res != addr {
		t.Error()
	}
}

func TestP2PKHAddress(t *testing.T) {

	h160, err := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	if err != nil {
		t.Error()
	}

	expected := "1BenRpVUFK65JFWcQSuHnJKzc4M8ZP8Eqa"
	addr := H160ToP2PKHAddress(h160, false)
	if addr != expected {
		t.Error()
	}

	expected = "mrAjisaT4LXL5MzE81sfcDYKU3wqWSvf9q"
	addr = H160ToP2PKHAddress(h160, true)
	if addr != expected {
		t.Error()
	}
}

func TestP2SHAddress(t *testing.T) {

	h160, err := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	if err != nil {
		t.Error()
	}

	expected := "3CLoMMyuoDQTPRD3XYZtCvgvkadrAdvdXh"
	addr := H160ToP2SHAddress(h160, false)
	if expected != addr {
		t.Error()
	}

	expected = "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7B"
	addr = H160ToP2SHAddress(h160, true)
	if expected != addr {
		t.Error()
	}
}
