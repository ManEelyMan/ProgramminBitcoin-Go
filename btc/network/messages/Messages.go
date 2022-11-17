package messages

import (
	"fmt"
	"io"
)

const PING_MESSAGE_NAME = "ping"
const PONG_MESSAGE_NAME = "pong"
const VERACK_MESSAGE_NAME = "verack"
const VERSION_MESSAGE_NAME = "version"
const GETHEADERS_MESSAGE_NAME = "getheaders"

type Message interface {
	GetName() string
	Serialize(io.Writer)
}

type MessageParser func(io.Reader) (Message, error)

var messageMap map[string]MessageParser

func init() {
	messageMap = make(map[string]MessageParser)
	messageMap[VERSION_MESSAGE_NAME] = parseVersionMessage
	messageMap[VERACK_MESSAGE_NAME] = parseVerackMessage
	messageMap[PING_MESSAGE_NAME] = parsePingMessage
	messageMap[PONG_MESSAGE_NAME] = parsePongMessage
	messageMap[GETHEADERS_MESSAGE_NAME] = parseGetHeadersMessage
}

func ParseMessagePayload(command string, reader io.Reader) (Message, error) {

	parser, ok := messageMap[command]
	if !ok {
		return nil, fmt.Errorf("command %v is unknown or unimplemented", command)
	}

	return parser(reader)
}
