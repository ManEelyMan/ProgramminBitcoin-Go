package ecc_test

import (
	"bitcoin-go/ecc"
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestPubPoint(t *testing.T) {

	point := N256P("5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", "6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da")
	secret := big.NewInt(7)
	mul := ecc.G.ScalarMultiply(secret)
	if !point.Equals(&mul) {
		t.Error()
	}

	point = N256P("c982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda", "7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55")
	secret = big.NewInt(1485)
	mul = ecc.G.ScalarMultiply(secret)
	if !point.Equals(&mul) {
		t.Error()
	}

	point = N256P("8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da", "662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82")
	secret = big.NewInt(2)
	secret = secret.Exp(secret, big.NewInt(128), nil)
	mul = ecc.G.ScalarMultiply(secret)
	if !point.Equals(&mul) {
		t.Error()
	}

	point = N256P("9577ff57c8234558f293df502ca4f09cbc65a6572c842b39b366f21717945116", "10b49c67fa9365ad7b90dab070be339a1daf9052373ec30ffae4f72d5e66d053")
	two := big.NewInt(2)
	firstExp := new(big.Int)
	secondExp := new(big.Int)
	firstExp = firstExp.Exp(two, big.NewInt(240), nil)
	secondExp = secondExp.Exp(two, big.NewInt(31), nil)
	secret = secret.Add(firstExp, secondExp)
	mul = ecc.G.ScalarMultiply(secret)
	if !point.Equals(&mul) {
		t.Error()
	}
}

func TestVerify(t *testing.T) {
	point := N256P("887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c", "61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	z := utility.HexStringToBigInt("ec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60")
	r := utility.HexStringToBigInt("ac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395")
	s := utility.HexStringToBigInt("68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4")
	sig := ecc.NewSignature(r, s)
	if !point.Verify(z, sig) {
		t.Error()
	}

	z = utility.HexStringToBigInt("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	r = utility.HexStringToBigInt("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c")
	s = utility.HexStringToBigInt("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6")
	sig = ecc.NewSignature(r, s)
	if !point.Verify(z, sig) {
		t.Error()
	}
}

func TestSEC(t *testing.T) {

	coef := new(big.Int)
	coef = coef.Exp(big.NewInt(999), big.NewInt(3), nil)
	if !SECTestCase(coef, "049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9", "039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5") {
		t.Error()
	}

	coef = big.NewInt(123)
	if !SECTestCase(coef, "04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b", "03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5") {
		t.Error()
	}

	coef = big.NewInt(42424242)
	if !SECTestCase(coef, "04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3", "03aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e") {
		t.Error()
	}
}

func TestSECParse(t *testing.T) {
	coef := new(big.Int)
	coef = coef.Exp(big.NewInt(999), big.NewInt(3), nil)

	if !SECParseTestCase(coef) {
		t.Error()
	}

	coef = big.NewInt(123)
	if !SECParseTestCase(coef) {
		t.Error()
	}

	coef = big.NewInt(42424242)
	if !SECParseTestCase(coef) {
		t.Error()
	}
}

func TestPointAddress(t *testing.T) {

	i := big.NewInt(888)
	i = i.Exp(i, big.NewInt(3), nil)
	if !PointAddressTestCase(i, "148dY81A9BmdpMhvYEVznrM45kWN32vSCN", "mieaqB68xDCtbUBYFoUNcmZNwk74xcBfTP", true) {
		t.Fail()
	}
	if !PointAddressTestCase(big.NewInt(321), "1S6g2xBJSED7Qr9CYZib5f4PYVhHZiVfj", "mfx3y63A7TfTtXKkv7Y6QzsPFY6QCBCXiP", false) {
		t.Fail()
	}
	if !PointAddressTestCase(big.NewInt(4242424242), "1226JSptcStqn4Yq9aAmNXdwdc2ixuH9nb", "mgY3bVusRUL6ZB2Ss999CSrGVbdRwVpM8s", false) {
		t.Fail()
	}
}

func N256P(xHex string, yHex string) ecc.Point {
	return ecc.NewSecp256k1Point(
		utility.HexStringToBigInt(xHex),
		utility.HexStringToBigInt(yHex))
}

func BytesFromHex(s string) []byte {
	a, err := hex.DecodeString(s)
	if err != nil {
		panic("Not a valid hex string!")
	}

	return a
}

func SECTestCase(coef *big.Int, expectedUncompressed string, expectedCompressed string) bool {

	fmt.Printf("Coef: %+v\n", coef)

	point := ecc.G.ScalarMultiply(coef)
	fmt.Printf("Point: %+v\n", point)

	uncompressed := BytesFromHex(expectedUncompressed)
	fmt.Printf("Uncomp expected: %+v\n", uncompressed)
	uncompSec := point.ToSEC(false)
	fmt.Printf("Uncomp actual: %+v\n", uncompSec)

	compressed := BytesFromHex(expectedCompressed)
	fmt.Printf("Comp expected: %+v\n", compressed)
	compSec := point.ToSEC(true)
	fmt.Printf("Comp actual: %+v\n", compSec)

	if !bytes.Equal(uncompressed, uncompSec) {
		return false
	}

	if !bytes.Equal(compressed, compSec) {
		return false
	}

	return true
}

func SECParseTestCase(coef *big.Int) bool {
	point := ecc.G.ScalarMultiply(coef)
	uncompressed := point.ToSEC(false)
	compressed := point.ToSEC(true)

	point2 := ecc.NewPointFromSEC(uncompressed)
	point3 := ecc.NewPointFromSEC(compressed)

	if !point.Equals(&point2) {
		return false
	}

	if !point.Equals(&point3) {
		return false
	}

	return true
}

func PointAddressTestCase(secret *big.Int, mainNetAddr string, testNetAddr string, compressed bool) bool {
	point := ecc.G.ScalarMultiply(secret)

	mainnet := point.Address(compressed, false)
	testnet := point.Address(compressed, true)

	if mainNetAddr != mainnet {
		return false
	}
	if testNetAddr != testnet {
		return false
	}

	return true
}
