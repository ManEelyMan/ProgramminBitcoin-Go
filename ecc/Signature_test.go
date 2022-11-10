package ecc_test

import (
	"bitcoin-go/ecc"
	"bitcoin-go/utility"
	"testing"
)

func TestDER(t *testing.T) {
	r := utility.HexStringToBigInt("ac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395")
	s := utility.HexStringToBigInt("68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4")
	sig := ecc.NewSignature(r, s)

	der := sig.ToDER()

	sig2 := ecc.NewSignatureFromDER(der)

	if !sig2.Equals(&sig) {
		t.Errorf("Didn't work")
	}
}
