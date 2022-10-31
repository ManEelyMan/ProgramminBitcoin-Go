package utility

import (
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
	// fmt.Printf("num %+v", num)

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
