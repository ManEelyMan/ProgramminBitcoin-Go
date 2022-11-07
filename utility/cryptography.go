package utility

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"math/big"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

func RandomData(numBytes int) []byte {
	bytes := make([]byte, numBytes)
	rand.Read(bytes)
	return bytes
}
func Sha1(bytes []byte) []byte {
	hash := sha1.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}
func Sha256(bytes []byte) []byte {
	hash := sha256.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}

func Hash256(bytes []byte) []byte {
	return Sha256(Sha256(bytes))
}

func HashRipemd160(bytes []byte) []byte {
	hash := ripemd160.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}

func Hash160(bytes []byte) []byte {
	return HashRipemd160(Sha256(bytes))
}

const BASE58_ALPHABET string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func DecodeBase58(addr string) ([]byte, bool) {

	var num *big.Int = big.NewInt(0)

	for _, c := range addr {
		var digit int = strings.IndexRune(BASE58_ALPHABET, c)
		if digit == -1 {
			return nil, false
		}

		num = num.Mul(num, big.NewInt(58))
		num = num.Add(num, big.NewInt(int64(digit)))
	}

	bin := num.Bytes()
	checksum := bin[len(bin)-4:]

	h256 := Hash256(bin[:len(bin)-4])

	if !bytes.Equal(checksum, h256[:4]) {
		return nil, false
	}

	return bin[1 : len(bin)-4], true
}

func EncodeBase58(bytes []byte) string {

	count := 0
	for _, b := range bytes {
		if b == 0x00 {
			count += 1
		} else {
			break
		}
	}

	num := new(big.Int)
	num.SetBytes(bytes)
	mod := big.NewInt(0)

	num.SetBytes(bytes)

	prefix := strings.Repeat("1", count)
	zero := big.NewInt(0)
	fiftyeight := big.NewInt(58)
	var result string = ""

	for num.Cmp(zero) > 0 {
		num, mod = num.DivMod(num, fiftyeight, mod)
		result = string(BASE58_ALPHABET[mod.Int64()]) + result
	}
	return prefix + result
}

func EncodeBase58Checksum(bytes []byte) string {
	concat := make([]byte, len(bytes)+4)
	copy(concat[0:], bytes)
	hash256 := Hash256(bytes)
	copy(concat[len(concat)-4:], hash256[:4])
	return EncodeBase58(concat)
}

func CreateBase58AddressFromHash(leadByte byte, hash []byte) string {
	concat := make([]byte, len(hash)+1)
	concat[0] = leadByte
	copy(concat[1:], hash)
	return EncodeBase58Checksum(concat)
}

func H160ToP2PKHAddress(hash []byte, testnet bool) string {
	return CreateBase58AddressFromHash(IIF(testnet, byte(0x6f), byte(0x00)).(byte), hash)
}

func H160ToP2SHAddress(hash []byte, testnet bool) string {
	return CreateBase58AddressFromHash(IIF(testnet, byte(0xc4), byte(0x05)).(byte), hash)
}
