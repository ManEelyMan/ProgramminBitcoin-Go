package messages

import "io"

type Ping struct {
	Nonce [8]byte
}

func (p Ping) GetName() string {
	return PING_MESSAGE_NAME
}

func (p Ping) Serialize(writer io.Writer) {
	writer.Write(p.Nonce[:])
}

func parsePingMessage(reader io.Reader) (Message, error) {
	ret := Ping{}
	_, err := reader.Read(ret.Nonce[:])
	return &ret, err
}
