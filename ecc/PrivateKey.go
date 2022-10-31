package ecc

import (
	"bitcoin-go/utility"
	"math/big"
)

type PrivateKey struct {
	Secret *big.Int
}

func NewPrivateKey(secret *big.Int) PrivateKey {
	return PrivateKey{Secret: secret}
}

func (key *PrivateKey) Sign(hash *big.Int) Signature {
	k := getK()
	r := (G.ScalarMultiply(k)).x
	k_inv := ModPowPrime(k, ModSubInt(N, 2), N)
	s := ModMulPrime(ModAdd(ModMul(r, key.Secret), hash), k_inv, N)

	half_n := ModMulInt(N, 2)
	if s.Cmp(half_n) > 0 {
		s = ModSub(N, s)
	}

	return Signature{R: r, S: s}
}

func (key *PrivateKey) WIF(compressed bool, testnet bool) string {

	bytes := make([]byte, 34)

	bytes[0] = utility.IIF(testnet, (byte)(0xef), (byte)(0x80)).(byte)

	copy(bytes[1:], key.Secret.Bytes())

	if compressed {
		bytes[33] = 0x01
	}

	length := len(bytes)
	if !compressed {
		length -= 1
	}

	return utility.EncodeBase58Checksum(bytes[:length])
}

func getK() *big.Int {

	r := utility.RandomData(32)
	i := new(big.Int)
	i.SetBytes(r)
	return i
}
