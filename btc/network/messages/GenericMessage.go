package messages

import "io"

// A generic message type for messages I haven't needed to fully implement yet.
type GenericMessage struct {
	name string
}

func (m GenericMessage) GetName() string     { return m.name }
func (m GenericMessage) Serialize(io.Writer) {}

func NewGenericMessage(command string) Message {
	return &GenericMessage{name: command}
}
