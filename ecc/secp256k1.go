package ecc

import (
	"bitcoin-go/utility"
	"math/big"
)

var A *big.Int
var B *big.Int
var P *big.Int
var Gx *big.Int
var Gy *big.Int
var G Point
var N *big.Int
var BigZero *big.Int
var BigOne *big.Int
var BigTwo *big.Int
var BigThree *big.Int

func init() {
	// These are used a lot, so why not treat them as wannabe "constants"
	BigZero = big.NewInt(0)
	BigOne = big.NewInt(1)
	BigTwo = big.NewInt(2)
	BigThree = big.NewInt(3)

	A = big.NewInt(0)
	B = big.NewInt(7)

	// Calculate our prime
	// 2 ^ 256 - 2 ^ 32 - 977
	tmp := big.NewInt(2)
	tmp = tmp.Exp(tmp, big.NewInt(256), nil)

	tmp2 := big.NewInt(2)
	tmp2 = tmp2.Exp(tmp2, big.NewInt(32), nil)

	P = tmp.Sub(tmp, tmp2)
	P = P.Sub(P, big.NewInt(977))

	Gx = utility.HexStringToBigInt("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	Gy = utility.HexStringToBigInt("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	G = NewSecp256k1Point(Gx, Gy)

	N = utility.HexStringToBigInt("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
}
