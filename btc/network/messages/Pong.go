package messages

import "io"

type Pong struct {
	Nonce [8]byte
}

func (p Pong) GetName() string {
	return PONG_MESSAGE_NAME
}

func (p Pong) Serialize(writer io.Writer) {
	writer.Write(p.Nonce[:])
}

func parsePongMessage(reader io.Reader) (Message, error) {
	ret := Pong{}
	_, err := reader.Read(ret.Nonce[:])
	return &ret, err
}
