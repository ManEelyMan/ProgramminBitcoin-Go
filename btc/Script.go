package btc

import (
	"bitcoin-go/utility"
	"io"
)

type Script struct {
	Data []byte
}

func ParseScript(reader io.Reader) Script {
	scriptLength := utility.ReadVarInt(reader)

	s := Script{}
	s.Data = make([]byte, scriptLength)
	reader.Read(s.Data)
	return s
}

func (script *Script) Serialize(writer io.Writer) {
	utility.WriteVarInt(writer, (uint64)(len(script.Data)))
	writer.Write(script.Data)
}
