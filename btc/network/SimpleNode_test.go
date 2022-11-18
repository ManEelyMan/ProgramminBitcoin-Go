package network

import "testing"

func TestSimpleNodeHandshake(t *testing.T) {
	node := NewSimpleNode("testnet.programmingbitcoin.com", 8333, true, false)

	node.Open()
	defer node.Close()

	err := node.Handshake()

	if err != nil {
		t.Error()
	}
}
