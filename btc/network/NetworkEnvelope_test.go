package network

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestParseNetworkEnvelope(t *testing.T) {

	msg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	reader := bytes.NewBuffer(msg)
	env, err := ParseNetworkEnvelope(reader)
	if err != nil {
		t.Error()
	}

	if env.Command != "verack" {
		t.Error()
	}

	if len(env.Payload) != 0 {
		t.Error()
	}

	msg, _ = hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	reader = bytes.NewBuffer(msg)
	env, err = ParseNetworkEnvelope(reader)
	if err != nil {
		t.Error()
	}

	if env.Command != "version" {
		t.Error()
	}

	if !bytes.Equal(msg[24:], env.Payload) {
		t.Error()
	}
}

func TestSerializeNetworkEnvelope(t *testing.T) {

	msg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	reader := bytes.NewBuffer(msg)
	env, err := ParseNetworkEnvelope(reader)
	if err != nil {
		t.Error()
	}

	res := env.Serialize()
	if !bytes.Equal(msg, res) {
		t.Error()
	}

	msg, _ = hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	reader = bytes.NewBuffer(msg)
	env, err = ParseNetworkEnvelope(reader)
	if err != nil {
		t.Error()
	}

	res = env.Serialize()
	if !bytes.Equal(msg, res) {
		t.Error()
	}
}
