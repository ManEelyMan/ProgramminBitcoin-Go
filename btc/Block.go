package btc

import (
	"bitcoin-go/utility"
	"bytes"
	"io"
)

type Block struct {
	Version       uint32
	PreviousBlock [32]byte
	MerkleRoot    [32]byte
	Timestamp     uint32
	Bits          uint32
	Nonce         uint32
}

func ParseBlock(reader io.Reader) (Block, error) {
	version := utility.ReadUint32(reader, true)
	prevBlock, err := utility.ReadBytes(reader, 32)
	if err != nil {
		return Block{}, err
	}
	prevBlock = utility.ReverseBytes(prevBlock)

	merkleRoot, err := utility.ReadBytes(reader, 32)
	if err != nil {
		return Block{}, err
	}
	merkleRoot = utility.ReverseBytes(merkleRoot)

	timestamp := utility.ReadUint32(reader, true)
	bits := utility.ReadUint32(reader, false)
	nonce := utility.ReadUint32(reader, false)

	block := Block{Version: version, Timestamp: timestamp, Bits: bits, Nonce: nonce}
	copy(block.PreviousBlock[:], prevBlock)
	copy(block.MerkleRoot[:], merkleRoot)

	return block, nil
}

func (b *Block) Serialize(writer io.Writer) {
	utility.WriteUint32(writer, b.Version, true)
	writer.Write(utility.ReverseBytes(b.PreviousBlock[:]))
	writer.Write(utility.ReverseBytes(b.MerkleRoot[:]))
	utility.WriteUint32(writer, b.Timestamp, true)
	utility.WriteUint32(writer, b.Bits, false)
	utility.WriteUint32(writer, b.Nonce, false)
}

func (b *Block) Hash() []byte {

	buffer := bytes.NewBuffer(make([]byte, 0))
	b.Serialize(buffer)

	hash := utility.ReverseBytes(utility.Hash256(buffer.Bytes()))
	return hash
}

func (b *Block) BIP09() bool {
	return b.Version>>29 == 0b001
}

func (b *Block) BIP91() bool {
	return (b.Version>>4)&1 == 1
}

func (b *Block) BIP141() bool {
	return (b.Version>>1)&1 == 1
}
