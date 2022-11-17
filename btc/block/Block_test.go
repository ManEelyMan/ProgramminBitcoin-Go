package block

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestBlockParse(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error()
	}
	buff := bytes.NewBuffer(rawBlock)
	block, err := ParseBlock(buff)
	if err != nil {
		t.Error()
	}

	if block.Version != 0x20000002 {
		t.Error()
	}

	expectedPrevBlock, err := hex.DecodeString("000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e")
	if err != nil {
		t.Error()
	}

	if !bytes.Equal(expectedPrevBlock, block.PreviousBlock[:]) {
		t.Error()
	}

	expectedMerkleRoot, err := hex.DecodeString("be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b")
	if err != nil {
		t.Error()
	}

	if !bytes.Equal(expectedMerkleRoot, block.MerkleRoot[:]) {
		t.Error()
	}

	if block.Timestamp != 0x59a7771e {
		t.Error()
	}

	expectedBits, _ := hex.DecodeString("e93c0118")
	if !bytes.Equal(block.Bits[:], expectedBits) {
		t.Error()
	}

	expectedNonce, _ := hex.DecodeString("a4ffd71d")
	if !bytes.Equal(block.Nonce[:], expectedNonce) {
		t.Error()
	}
}

func TestBlockSerialize(t *testing.T) {

	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error()
	}

	reader := bytes.NewBuffer(rawBlock)
	block, err := ParseBlock(reader)
	if err != nil {
		t.Error()
	}

	writer := bytes.NewBuffer(make([]byte, 0))
	block.Serialize(writer)
	ser := writer.Bytes()

	if !bytes.Equal(ser, rawBlock) {
		t.Error()
	}
}

func TestHash(t *testing.T) {

	block := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	hash := block.Hash()
	expected, err := hex.DecodeString("0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523")
	if err != nil {
		t.Error()
	}

	if !bytes.Equal(expected, hash) {
		t.Error()
	}
}

func TestBlockBIP09(t *testing.T) {
	block1 := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if !block1.BIP09() {
		t.Error()
	}

	block2 := blockFromHexString("0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27")
	if block2.BIP09() {
		t.Error()
	}
}

func TestBlockBIP91(t *testing.T) {

	block1 := blockFromHexString("1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1")
	if !block1.BIP91() {
		t.Error()
	}

	block2 := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if block2.BIP91() {
		t.Error()
	}
}

func TestBlockBIP141(t *testing.T) {

	block1 := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if !block1.BIP141() {
		t.Error()
	}

	block2 := blockFromHexString("0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618")
	if block2.BIP141() {
		t.Error()
	}
}

func TestBlockCalculateNewBits(t *testing.T) {

	var prevBits [4]byte
	bits, _ := hex.DecodeString("54d80118")
	copy(prevBits[:], bits)

	newBits := CalculateNewBits(prevBits, 302400)

	expectedBits, _ := hex.DecodeString("00157617")

	if !bytes.Equal(expectedBits, newBits[:]) {
		t.Error()
	}
}

func TestBlockTarget(t *testing.T) {
	block := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	target := block.Target()

	expected := utility.HexStringToBigInt("13ce9000000000000000000000000000000000000000000")

	if expected.Cmp(target) != 0 {
		t.Error()
	}
}

func TestBlockDifficulty(t *testing.T) {
	block := blockFromHexString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	diff := big.NewInt(888171856257)
	if block.Difficulty().Cmp(diff) != 0 {
		t.Error()
	}
}

func TestBlockCheckProofOfWork(t *testing.T) {
	block := blockFromHexString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1")
	if !block.CheckProofOfWork() {
		t.Error()
	}

	block = blockFromHexString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0")
	if block.CheckProofOfWork() {
		t.Error()
	}
}

func blockFromHexString(s string) Block {
	raw, _ := hex.DecodeString(s)
	reader := bytes.NewBuffer(raw)
	block, _ := ParseBlock(reader)
	return block
}
