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

	merkleRoot, err := utility.ReadBytes(reader, 32)
	if err != nil {
		return Block{}, err
	}

	timestamp := utility.ReadUint32(reader, true)
	bits := utility.ReadUint32(reader, true)
	nonce := utility.ReadUint32(reader, true)

	block := Block{Version: version, Timestamp: timestamp, Bits: bits, Nonce: nonce}
	copy(block.PreviousBlock[:], prevBlock)
	copy(block.MerkleRoot[:], merkleRoot)

	return block, nil
}

func (b *Block) Serialize(writer io.Writer) {
	utility.WriteUint32(writer, b.Version, true)
	writer.Write(b.PreviousBlock[:])
	writer.Write(b.MerkleRoot[:])
	utility.WriteUint32(writer, b.Timestamp, true)
	utility.WriteUint32(writer, b.Bits, true)
	utility.WriteUint32(writer, b.Nonce, true)
}

func (b *Block) Hash() []byte {

	buffer := bytes.NewBuffer(make([]byte, 0))
	b.Serialize(buffer)

	hash := utility.Hash256(utility.ReverseBytes(buffer.Bytes()))
	return hash
}
