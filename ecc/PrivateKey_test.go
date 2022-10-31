package ecc_test

import (
	"bitcoin-go/ecc"
	"bitcoin-go/utility"
	"math/big"
	"testing"
)

func TestPrivateKeyWIF(t *testing.T) {

	two := big.NewInt(2)
	first := new(big.Int)
	first = first.Exp(two, big.NewInt(256), nil)
	second := new(big.Int)
	second = second.Exp(two, big.NewInt(199), nil)

	secret := new(big.Int)
	secret = secret.Sub(first, second)
	if !WIFTestCase(secret, true, false, "L5oLkpV3aqBJ4BgssVAsax1iRa77G5CVYnv9adQ6Z87te7TyUdSC") {
		t.Error()
	}

	second = second.Exp(two, big.NewInt(201), nil)
	secret = secret.Sub(first, second)
	if !WIFTestCase(secret, false, true, "93XfLeifX7Jx7n7ELGMAf1SUR6f9kgQs8Xke8WStMwUtrDucMzn") {
		t.Error()
	}

	secret = utility.HexStringToBigInt("1cca23de92fd1862fb5b76e5f4f50eb082165e5191e116c18ed1a6b24be6a53f")
	if !WIFTestCase(secret, true, true, "cNYfWuhDpbNM1JWc3c6JTrtrFVxU4AGhUKgw5f93NP2QaBqmxKkg") {
		t.Error()
	}
}

func WIFTestCase(secret *big.Int, compressed bool, testnet bool, expectedWif string) bool {

	pk := ecc.NewPrivateKey(secret)
	wif := pk.WIF(compressed, testnet)

	return wif == expectedWif
}
