package messages

import (
	"bitcoin-go/utility"
	"io"
)

type Version struct {
	Version                 uint32
	SenderNetworkServices   uint64
	Timestamp               uint64
	ReceiverNetworkServices uint64
	ReceiverNetworkAddress  [16]byte
	ReceiverNetworkPort     [2]byte
	SenderNetworkServices2  uint64
	SenderNetworkAddress    [16]byte
	SenderNetworkPort       [2]byte
	Nonce                   [8]byte
	UserAgent               string
	Height                  [4]byte
	RelayFlag               byte
}

func parseVersionMessage(reader io.Reader) (Message, error) {

	msg := new(Version)
	msg.Version = utility.ReadUint32(reader, true)
	msg.SenderNetworkServices = utility.ReadUint64(reader, true)
	msg.Timestamp = utility.ReadUint64(reader, true)
	msg.ReceiverNetworkServices = utility.ReadUint64(reader, true)
	_, err := reader.Read(msg.ReceiverNetworkAddress[:])
	if err != nil {
		return msg, err
	}
	_, err = reader.Read(msg.ReceiverNetworkPort[:])
	if err != nil {
		return msg, err
	}
	msg.SenderNetworkServices2 = utility.ReadUint64(reader, true)
	_, err = reader.Read(msg.SenderNetworkAddress[:])
	if err != nil {
		return msg, err
	}
	_, err = reader.Read(msg.SenderNetworkPort[:])
	if err != nil {
		return msg, err
	}
	_, err = reader.Read(msg.Nonce[:])
	if err != nil {
		return msg, err
	}
	ua, err := utility.ReadBytes(reader, 27)
	if err != nil {
		return msg, err
	}
	msg.UserAgent = string(ua)
	_, err = reader.Read(msg.Height[:])
	if err != nil {
		return msg, err
	}
	msg.RelayFlag = utility.ReadByte(reader)

	return msg, nil
}

func (m Version) GetName() string {
	return VERSION_MESSAGE_NAME
}

func (m Version) Serialize(writer io.Writer) {
	utility.WriteUint32(writer, m.Version, true)
	utility.WriteUint64(writer, m.SenderNetworkServices, true)
	utility.WriteUint64(writer, m.Timestamp, true)
	utility.WriteUint64(writer, m.ReceiverNetworkServices, true)
	_, _ = writer.Write(m.ReceiverNetworkAddress[:])
	_, _ = writer.Write(m.ReceiverNetworkPort[:])
	utility.WriteUint64(writer, m.SenderNetworkServices2, true)
	_, _ = writer.Write(m.SenderNetworkAddress[:])
	_, _ = writer.Write(m.SenderNetworkPort[:])
	_, _ = writer.Write(m.Nonce[:])
	utility.WriteVarInt(writer, uint64(len(m.UserAgent)))
	_, _ = writer.Write([]byte(m.UserAgent))
	_, _ = writer.Write(m.Height[:])
	tmp := make([]byte, 1)
	tmp[0] = m.RelayFlag
	_, _ = writer.Write(tmp)
}
