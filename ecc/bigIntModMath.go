package ecc

import (
	"math/big"
)

func ModAdd(x *big.Int, y *big.Int) *big.Int {
	z := new(big.Int)
	z = z.Add(x, y)
	return z.Mod(z, P)
}

func ModAddInt(x *big.Int, y int64) *big.Int {
	return ModAdd(x, big.NewInt(y))
}

func ModSub(x *big.Int, y *big.Int) *big.Int {
	z := new(big.Int)
	z = z.Sub(x, y)
	return z.Mod(z, P)
}

func ModSubInt(x *big.Int, y int64) *big.Int {
	return ModSub(x, big.NewInt(y))
}

func ModMulPrime(x *big.Int, y *big.Int, prime *big.Int) *big.Int {
	z := new(big.Int)
	z.Mul(x, y)
	return z.Mod(z, prime)
}

func ModMul(x *big.Int, y *big.Int) *big.Int {
	return ModMulPrime(x, y, P)
}

func ModMulInt(x *big.Int, y int64) *big.Int {
	return ModMul(x, big.NewInt(y))
}

func ModPowPrime(x *big.Int, pow *big.Int, prime *big.Int) *big.Int {
	z := new(big.Int)
	return z.Exp(x, pow, prime)
}

func ModPow(x *big.Int, pow *big.Int) *big.Int {
	return ModPowPrime(x, pow, P)
}

func ModPowInt(x *big.Int, pow int64) *big.Int {
	return ModPow(x, big.NewInt(pow))
}

func ModPowPrimeInt(x *big.Int, pow int64, prime *big.Int) *big.Int {
	return ModPowPrime(x, big.NewInt(pow), prime)
}

func ModDiv(x *big.Int, y *big.Int) *big.Int {
	power := ModSubInt(P, 2)
	inverse := ModPowPrime(y, power, P)
	return ModMul(x, inverse)
}

func ModDivInt(x *big.Int, y int64) *big.Int {
	return ModDiv(x, big.NewInt(y))
}

func ModSqrt(x *big.Int) *big.Int {
	exp := ModDivInt(ModAddInt(P, 1), 4)
	return ModPow(x, exp)
}

func IsEven(x *big.Int) bool {
	mod := new(big.Int)
	mod = mod.Mod(x, big.NewInt(2))
	isEven := mod.Cmp(big.NewInt(0)) == 0
	return isEven
}
