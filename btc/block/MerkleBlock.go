package block

import (
	"bitcoin-go/utility"
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

func (mb *MerkleBlock) IsValid() bool {

	panic("not yet implemented")

	/*      '''Verifies whether the merkle tree information validates to the merkle root'''
	        # convert the flags field to a bit field
	        flag_bits = bytes_to_bit_field(self.flags)
	        # reverse self.hashes for the merkle root calculation
	        hashes = [h[::-1] for h in self.hashes]
	        # initialize the merkle tree
	        merkle_tree = MerkleTree(self.total)
	        # populate the tree with flag bits and hashes
	        merkle_tree.populate_tree(flag_bits, hashes)
	        # check if the computed root reversed is the same as the merkle root
	        return merkle_tree.root()[::-1] == self.merkle_root
	*/
}
