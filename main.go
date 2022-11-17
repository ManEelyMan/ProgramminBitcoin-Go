package main

import (
	"bitcoin-go/btc/network"
	"bitcoin-go/btc/network/messages"
	"fmt"
	"io"
)

func main() {

	node := network.NewSimpleNode("testnet.programmingbitcoin.com", 0, true, false)
	err := node.Open()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	defer node.Close()

	versionMsg := messages.Version{}
	node.Send(versionMsg)

	verackReceived := false
	versionReceived := false

	for !verackReceived && !versionReceived {
		message, err := node.WaitFor([]string{messages.VERACK_MESSAGE_NAME, messages.VERSION_MESSAGE_NAME})
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached end of stream. No data.")
			} else {
				fmt.Printf("%v\n", message)
			}
			break
		}

		if message.GetName() == messages.VERACK_MESSAGE_NAME {
			verackReceived = true
		} else if message.GetName() == messages.VERSION_MESSAGE_NAME {
			versionReceived = true
		}
	}

	fmt.Println("Done!")
}
