package network

import (
	"bitcoin-go/btc/network/messages"
	"bytes"
	"fmt"
	"net"
)

type SimpleNode struct {
	Host       string
	Port       uint16
	TestNet    bool
	Logging    bool
	connection net.Conn
}

func NewSimpleNode(host string, port uint16, testnet bool, logging bool) *SimpleNode {
	newNode := new(SimpleNode)
	newNode.Host = host

	if port == 0 {
		if testnet {
			newNode.Port = 18333
		} else {
			newNode.Port = 8333
		}
	} else {
		newNode.Port = port
	}

	newNode.TestNet = testnet
	newNode.Logging = logging

	return newNode
}

func (n *SimpleNode) Open() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", n.Host, n.Port))
	if err == nil {
		n.connection = conn
	}
	return err
}

func (n *SimpleNode) Close() error {
	return n.connection.Close()
}

func (n *SimpleNode) Send(message messages.Message) error {
	writer := bytes.NewBuffer(make([]byte, 0))
	message.Serialize(writer)

	env := NewNetworkEnvelope(message.GetName(), writer.Bytes(), n.TestNet)

	_, err := n.connection.Write(env.Serialize())
	return err
}

func (n *SimpleNode) Read() (*NetworkEnvelope, error) {

	readBuffer := make([]byte, 2048) // Twice the maximum block size.  TODO: Make this static for reuse (unless I multithread this)
	count, err := n.connection.Read(readBuffer)
	if err != nil {
		return nil, err
	}

	if count >= len(readBuffer) {
		panic("We maxed out our read buffer!")
	}

	reader := bytes.NewBuffer(readBuffer)
	return ParseNetworkEnvelope(reader)
}

func (n *SimpleNode) WaitFor(messageTypes []string) (messages.Message, error) {

	for { // Loop until we get one of the messages we want (or forever!!)
		env, err := n.Read()
		if err != nil {
			return nil, err
		}

		msgName := env.GetMessage().GetName()

		switch msgName {

		case messages.VERSION_MESSAGE_NAME:
			n.Send(messages.Verack{})

		case messages.PING_MESSAGE_NAME:
			n.Send(messages.Pong{Nonce: (env.GetMessage().(messages.Ping)).Nonce}) // I think I'm supposed to reply with the same Nonce as the ping?

		default:
			idx := findStringInSlice(msgName, messageTypes)
			if idx > -1 {
				return env.GetMessage(), nil
			}
		}
	}
}

func findStringInSlice(s string, arr []string) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == s {
			return i
		}
	}
	return -1
}
