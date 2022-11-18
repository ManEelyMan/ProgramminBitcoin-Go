package network

import (
	"bitcoin-go/btc/network/messages"
	"bitcoin-go/utility"
	"bytes"
	"errors"
	"io"
	"strings"
)

var MainnetMagic []byte
var TestnetMagic []byte

func init() {
	MainnetMagic = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	TestnetMagic = []byte{0x0b, 0x11, 0x09, 0x07}
}

type NetworkEnvelope struct {
	Magic          []byte
	Command        string
	Payload        []byte
	payloadMessage *messages.Message
}

func NewNetworkEnvelope(command string, payload []byte, testnet bool) NetworkEnvelope {
	magic := utility.IIF(testnet, TestnetMagic, MainnetMagic).([]byte)
	return NetworkEnvelope{Command: command, Payload: payload, Magic: magic}
}

func ParseNetworkEnvelope(reader io.Reader) (*NetworkEnvelope, error) {
	magic, err := utility.ReadBytes(reader, 4)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(MainnetMagic, magic) && !bytes.Equal(TestnetMagic, magic) {
		return nil, errors.New("magic value matches neither mainnet nor testnet")
	}

	command, err := utility.ReadBytes(reader, 12)
	if err != nil {
		return nil, err
	}

	len := utility.ReadUint32(reader, true)

	checksum, err := utility.ReadBytes(reader, 4)
	if err != nil {
		return nil, err
	}

	payload, err := utility.ReadBytes(reader, uint(len))
	if err != nil {
		return nil, err
	}

	hash := utility.Hash256(payload)

	if !bytes.Equal(checksum, hash[:4]) {
		return nil, errors.New("checksum mismatch for payload")
	}

	ne := new(NetworkEnvelope)
	ne.Magic = magic
	ne.Command = strings.TrimRight(string(command), "\000")
	ne.Payload = payload

	return ne, nil
}

func (env *NetworkEnvelope) Serialize() []byte {
	writer := bytes.NewBuffer(make([]byte, 0))
	writer.Write(env.Magic)

	commandLength := len(env.Command)
	tmp := env.Command + strings.Repeat("\000", 12-commandLength)
	writer.WriteString(tmp)

	payloadLen := len(env.Payload)
	payloadHash := utility.Hash256(env.Payload)

	utility.WriteUint32(writer, uint32(payloadLen), true)
	writer.Write(payloadHash[:4])
	writer.Write(env.Payload)

	return writer.Bytes()
}

func (env *NetworkEnvelope) GetMessage() messages.Message {
	if env.payloadMessage != nil {
		msg, _ := messages.ParseMessagePayload(env.Command, bytes.NewBuffer(env.Payload))
		env.payloadMessage = &msg
	} else {
		msg := messages.NewGenericMessage(env.Command)
		env.payloadMessage = &msg
	}

	return *env.payloadMessage
}
