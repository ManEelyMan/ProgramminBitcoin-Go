package utility

import (
	"encoding/binary"
	"io"
	"math"
	"math/big"
)

func ReadVarInt(reader io.Reader) uint64 {
	intLen := 1

	// Determine length of the size in bytes.
	buffer := make([]byte, 8)
	reader.Read(buffer[:1])

	if (buffer[0]) >= 253 {
		switch buffer[0] {
		case 0xff:
			intLen = 8
		case 0xfe:
			intLen = 4
		case 0xfd:
			intLen = 2
		}

		// Read in the bytes for the number.
		reader.Read(buffer[0:intLen])
	}

	varInt := binary.LittleEndian.Uint64(buffer)
	return varInt
}

func ReadBigInt(reader io.Reader, littleEndian bool) *big.Int {

	buffer := make([]byte, 32)
	reader.Read(buffer)

	if littleEndian {
		ReverseBytes(buffer)
	}

	i := new(big.Int)
	i.SetBytes(buffer)

	return i
}

func ReadByte(reader io.Reader) byte {
	buff := make([]byte, 1)
	reader.Read(buff)
	return buff[0]
}

func ReadBytes(reader io.Reader, len uint) ([]byte, error) {
	buff := make([]byte, len)
	_, e := reader.Read(buff)
	return buff, e
}

func ReadUInt16(reader io.Reader, littleEndian bool) uint16 {
	buffer := make([]byte, 2)
	reader.Read(buffer)

	if littleEndian {
		return binary.LittleEndian.Uint16(buffer)
	} else {
		return binary.BigEndian.Uint16(buffer)
	}
}

func ReadUint32(reader io.Reader, littleEndian bool) uint32 {

	buffer := make([]byte, 4)
	reader.Read(buffer)

	if littleEndian {
		return binary.LittleEndian.Uint32(buffer)
	} else {
		return binary.BigEndian.Uint32(buffer)
	}
}

func ReadUint64(reader io.Reader, littleEndian bool) uint64 {
	buffer := make([]byte, 8)
	reader.Read(buffer)

	if littleEndian {
		return binary.LittleEndian.Uint64(buffer)
	} else {
		return binary.BigEndian.Uint64(buffer)
	}
}

func WriteVarInt(writer io.Writer, i uint64) {

	buff := make([]byte, 9)
	binary.LittleEndian.PutUint64(buff[1:], i)

	varIntLen := 0
	buff[0] = buff[1]

	if i > 253 {
		if i > uint64(math.Pow(2, 32)) {
			varIntLen = 8
			buff[0] = 0xff
		} else if i > uint64(math.Pow(2, 16)) {
			varIntLen = 4
			buff[0] = 0xfe
		} else {
			varIntLen = 2
			buff[0] = 0xfd
		}
	}

	writer.Write(buff[:varIntLen+1])
}

func WriteBigInt(writer io.Writer, i *big.Int, littleEndian bool) {
	buff := i.Bytes()

	if littleEndian {
		ReverseBytes(buff)
	}

	writer.Write(buff)
}

func WriteUint32(writer io.Writer, i uint32, littleEndian bool) {
	buff := make([]byte, 4)

	if littleEndian {
		binary.LittleEndian.PutUint32(buff, i)
	} else {
		binary.BigEndian.PutUint32(buff, i)
	}

	writer.Write(buff)
}

func WriteUint64(writer io.Writer, i uint64, littleEndian bool) {
	buff := make([]byte, 8)

	if littleEndian {
		binary.LittleEndian.PutUint64(buff, i)
	} else {
		binary.BigEndian.PutUint64(buff, i)
	}

	writer.Write(buff)
}
