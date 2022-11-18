package messages

import (
	"bitcoin-go/btc/block"
	"bitcoin-go/utility"
	"fmt"
	"io"
)

type Headers struct {
	Blocks []block.Block
}

func parseHeadersMessage(reader io.Reader) (Message, error) {

	numHeaders := utility.ReadVarInt(reader)
	blocks := make([]block.Block, numHeaders)

	for i := 0; i < int(numHeaders); i++ {
		block, err := block.ParseBlock(reader)
		if err != nil {
			return nil, err
		}

		blocks[i] = block

		numTxs := utility.ReadVarInt(reader)
		if numTxs != 0 {
			return nil, fmt.Errorf("the number of transactions should be zero")
		}
	}

	return Headers{Blocks: blocks}, nil
}

func (h Headers) GetName() string {
	return HEADERS_MESSAGE_NAME
}

func (h Headers) Serialize(writer io.Writer) {
	// Not needed.
	panic("headers serialization not implemented")
}
