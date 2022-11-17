package messages

import (
	"bitcoin-go/utility"
	"fmt"
	"io"
)

type GetHeaders struct {
	Version    uint32
	NumHashes  uint32
	StartBlock [32]byte
	EndBlock   [32]byte
}

func parseGetHeadersMessage(reader io.Reader) (Message, error) {
	// TODO: I don't think this is something we ever have to parse
	return nil, fmt.Errorf("not implmented")
}

func (gh GetHeaders) GetName() string {
	return GETHEADERS_MESSAGE_NAME
}

func (gh GetHeaders) Serialize(writer io.Writer) {
	utility.WriteUint32(writer, gh.Version, true)
	utility.WriteVarInt(writer, uint64(gh.NumHashes))

	// Have to reverse the bytes.
	var dupe [32]byte
	copy(dupe[:], gh.StartBlock[:])
	utility.ReverseBytes(dupe[:])
	writer.Write(dupe[:])

	if len(gh.EndBlock) == 0 {
		writer.Write(make([]byte, 32))
	} else {
		//Reverse these too!
		copy(dupe[:], gh.EndBlock[:])
		utility.ReverseBytes(dupe[:])
		writer.Write(dupe[:])
	}
}
