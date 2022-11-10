package utility

import (
	"math/big"
)

func HexStringToBigInt(s string) *big.Int {
	tmp := new(big.Int)
	tmp, ok := tmp.SetString(s, 16)

	if !ok {
		panic("Set string error!")
	}

	return tmp
}

func ReverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}

	return b
}

func IIF(condition bool, ifTrue interface{}, ifFalse interface{}) interface{} {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}
