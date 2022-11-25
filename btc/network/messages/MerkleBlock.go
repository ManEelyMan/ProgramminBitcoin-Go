package messages

import (
	"bitcoin-go/utility"
	"bytes"
	"io"
)

type MerkleBlock struct {
	Version       uint32
	PreviousBlock [32]byte
	MerkleRoot    [32]byte
	Timestamp     uint32
	Bits          [4]byte
	Nonce         [4]byte
	Total         uint32
	Hashes        [][32]byte
	Flags         []byte
	isValid       *bool
}

func ParseMerkleBlock(reader io.Reader) (MerkleBlock, error) {

	block := MerkleBlock{}

	// Version
	block.Version = utility.ReadUint32(reader, true)

	// PreviousBlock
	_, err := reader.Read(block.PreviousBlock[:])
	if err != nil {
		return block, err
	}
	utility.ReverseBytes(block.PreviousBlock[:])

	// MerkleRoot
	_, err = reader.Read(block.MerkleRoot[:])
	if err != nil {
		return block, err
	}
	utility.ReverseBytes(block.MerkleRoot[:])

	// Timestamp
	block.Timestamp = utility.ReadUint32(reader, true)

	// Bits
	_, err = reader.Read(block.Bits[:])
	if err != nil {
		return block, err
	}

	// Nonce
	_, err = reader.Read(block.Nonce[:])
	if err != nil {
		return block, err
	}

	// Total
	block.Total = utility.ReadUint32(reader, true)

	// Hashes
	numHashes := utility.ReadVarInt(reader)
	block.Hashes = make([][32]byte, numHashes)

	for i := 0; i < int(numHashes); i++ {
		_, err = reader.Read(block.Hashes[i][:])
		if err != nil {
			return block, err
		}
		// They're little-endian, so reverse them.
		utility.ReverseBytes(block.Hashes[i][:])
	}

	// Flags
	flagLength := utility.ReadVarInt(reader)
	block.Flags, err = utility.ReadBytes(reader, uint(flagLength))
	if err != nil {
		return block, err
	}

	return block, nil
}

func parseMerkleBlockMessage(reader io.Reader) (Message, error) {
	return ParseMerkleBlock(reader)
}

func (mb MerkleBlock) GetName() string {
	return MERKLE_BLOCK_MESSAGE_NAME
}
func (mb MerkleBlock) Serialize(io.Writer) {
	panic("serialize not implemented")
}

func (mb *MerkleBlock) IsValid() bool {

	// Return cached value
	if mb.isValid != nil {
		return *mb.isValid
	}

	// Reverse the hashes
	for i := 0; i < len(mb.Hashes); i++ {
		utility.ReverseBytes(mb.Hashes[i][:])
	}

	tree := NewMerkleTree(mb.Total)
	tree.PopulateTree(utility.MerkelBytesToBitField(mb.Flags), mb.Hashes)

	root := tree.Root()
	derefRoot := *root
	utility.ReverseBytes(derefRoot[:])

	isValid := bytes.Equal(derefRoot[:], mb.MerkleRoot[:])
	mb.isValid = &isValid
	return *mb.isValid
}
