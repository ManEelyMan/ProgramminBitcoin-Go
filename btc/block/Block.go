package block

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/binary"
	"io"
	"math/big"
)

type Block struct {
	Version       uint32
	PreviousBlock [32]byte
	MerkleRoot    [32]byte
	Timestamp     uint32
	Bits          [4]byte
	Nonce         [4]byte
}

var TWO_WEEKS uint
var EIGHT_WEEKS uint
var HALF_WEEK uint
var MAX_TARGET *big.Int

func init() {
	tmp := big.NewInt(256)
	tmp = tmp.Exp(tmp, big.NewInt(0x1d-0x03), nil)
	MAX_TARGET = tmp.Mul(big.NewInt(0xffff), tmp)

	TWO_WEEKS = 60 * 60 * 24 * 14
	EIGHT_WEEKS = TWO_WEEKS * 4
	HALF_WEEK = TWO_WEEKS / 4
}

func ParseBlock(reader io.Reader) (Block, error) {
	block := Block{}

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

	return block, nil
}

func (b *Block) Serialize(writer io.Writer) {
	utility.WriteUint32(writer, b.Version, true)
	writer.Write(utility.ReverseBytes(b.PreviousBlock[:]))
	writer.Write(utility.ReverseBytes(b.MerkleRoot[:]))
	utility.WriteUint32(writer, b.Timestamp, true)
	writer.Write(b.Bits[:])
	writer.Write(b.Nonce[:])
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

func (b *Block) Target() *big.Int {
	return bits2Target(b.Bits)
}

func (b *Block) Difficulty() *big.Int {
	return MAX_TARGET.Div(MAX_TARGET, b.Target())
}

func (b *Block) CheckProofOfWork() bool {

	// Serialize
	buffer := bytes.NewBuffer(make([]byte, 0))
	b.Serialize(buffer)
	byteArr := buffer.Bytes()

	// Hash
	byteArr = utility.Hash256(byteArr)

	// Convert to int
	buffer = bytes.NewBuffer(byteArr)
	proof := utility.ReadBigInt(buffer, true)

	// Compare to our target.
	return proof.Cmp(b.Target()) < 0
}

func bits2Target(bits [4]byte) *big.Int {
	exponent := bits[3] - 3
	bits[3] = 0
	coeff := binary.LittleEndian.Uint32(bits[:])

	tmp := big.NewInt(int64(exponent))
	tmp = tmp.Exp(big.NewInt(256), tmp, nil)

	result := tmp.Mul(big.NewInt(int64(coeff)), tmp)
	return result
}

func target2Bits(target *big.Int) [4]byte {

	targetBytes := target.Bytes()
	var newBits [4]byte

	if targetBytes[0] > 0x7f {
		newBits[0] = targetBytes[1]
		newBits[1] = targetBytes[0]
		newBits[2] = 0x00
		newBits[3] = byte(len(targetBytes)) + 1
	} else {
		newBits[0] = targetBytes[2]
		newBits[1] = targetBytes[1]
		newBits[2] = targetBytes[0]
		newBits[3] = byte(len(targetBytes))
	}
	return newBits
}

func CalculateNewBits(prevBits [4]byte, timeDiff uint) [4]byte {

	// Cap the timeDiff value to between HALF_WEEK to EIGHT_WEEKS
	if timeDiff > EIGHT_WEEKS {
		timeDiff = EIGHT_WEEKS
	} else if timeDiff < HALF_WEEK {
		timeDiff = HALF_WEEK
	}

	tmp := bits2Target(prevBits)
	tmp = tmp.Mul(tmp, big.NewInt(int64(timeDiff)))
	newTarget := tmp.Div(tmp, big.NewInt(int64(TWO_WEEKS)))

	if newTarget.Cmp(MAX_TARGET) > 0 {
		newTarget = MAX_TARGET
	}

	return target2Bits(newTarget)
}
