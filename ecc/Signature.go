package ecc

import (
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func NewSignature(r *big.Int, s *big.Int) Signature {
	return Signature{R: r, S: s}
}

func (sig *Signature) Equals(sig2 *Signature) bool {
	if sig2 == nil {
		return false
	}
	if sig.R.Cmp(sig2.R) != 0 {
		return false
	}
	if sig.S.Cmp(sig2.S) != 0 {
		return false
	}
	return true
}

func (sig *Signature) ToDER() []byte {

	var buffer []byte

	rbytes, rzeros, rsign := analyzeBigInt(sig.R)
	rauglen := len(rbytes) - rzeros + rsign
	sbytes, szeros, ssign := analyzeBigInt(sig.S)
	sauglen := len(sbytes) - szeros + ssign

	// Create the buffer the correct size. It has three markers, three 1-byte lengths
	// and then the augmented sizes of the R and S values
	buffer = make([]byte, 3+3+rauglen+sauglen)

	// Fill in our buffer
	buffer[0] = 0x30 // Marker
	buffer[1] = byte(rauglen + sauglen + 4)

	// Fill in the R value.
	index := 2 + fillInValue(buffer[2:], rbytes[rzeros:], rsign)

	// Fill in the S value.
	fillInValue(buffer[index:], sbytes[szeros:], ssign)

	return buffer
}

func NewSignatureFromDER(bytes []byte) Signature {

	// TODO: Fix this to return either (Signature, error) or (Signature, bool) to handle errors gracefully.
	if bytes[0] != 0x30 {
		panic("Not a DER format")
	}

	startOfR := 2
	startOfS := 4 + (int)(bytes[3])

	r := extractBigIntFromBytes(bytes[startOfR:])
	s := extractBigIntFromBytes(bytes[startOfS:])

	return NewSignature(r, s)
}

func analyzeBigInt(i *big.Int) ([]byte, int, int) {

	// Get the bytes
	var buffer []byte = i.Bytes()

	// Track the leading zeros
	zeros := 0
	for zeros = 0; buffer[zeros] == 0x00; zeros++ {
	}

	// See if we need a byte for a negative number
	signByte := 0
	if i.Sign() == -1 {
		signByte = 1
	}

	return buffer, zeros, signByte
}

func fillInValue(dest []byte, bytes []byte, sign int) int {

	dest[0] = 0x02
	dest[1] = byte(len(bytes))
	if sign > 0 {
		dest[2] = 0x00
	}
	copy(dest[2+sign:], bytes)

	return len(bytes) + 2 + sign
}

func extractBigIntFromBytes(buffer []byte) *big.Int {

	if buffer[0] != 0x02 {
		panic("Not a valid DER encoding.")
	}

	len := (int)(buffer[1])
	intOffset := 2

	// Account for the padding of a negative number
	if buffer[intOffset] == 0x00 {
		intOffset += 1
		len -= 1
	}

	i := new(big.Int)
	i.SetBytes(buffer[intOffset : intOffset+len])

	return i
}
