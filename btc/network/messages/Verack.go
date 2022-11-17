package messages

import "io"

type Verack struct {
	// Empty by design.
}

func parseVerackMessage(reader io.Reader) (Message, error) {
	return &Verack{}, nil
}

func (m Verack) GetName() string {
	return VERACK_MESSAGE_NAME
}

func (m Verack) Serialize(_ io.Writer) {
}
