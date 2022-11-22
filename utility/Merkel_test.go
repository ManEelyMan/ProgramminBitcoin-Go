package utility

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMerkleParent(t *testing.T) {

	left, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	right, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	expected, _ := hex.DecodeString("8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd")

	parent := MerkelParent(left, right)

	if !bytes.Equal(expected, parent) {
		t.Error()
	}
}

func TestMerkleParentLevel(t *testing.T) {

	var hashes [][]byte = make([][]byte, 11)
	hashes[0], _ = hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hashes[1], _ = hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	hashes[2], _ = hex.DecodeString("f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0")
	hashes[3], _ = hex.DecodeString("3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181")
	hashes[4], _ = hex.DecodeString("10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae")
	hashes[5], _ = hex.DecodeString("7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161")
	hashes[6], _ = hex.DecodeString("8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc")
	hashes[7], _ = hex.DecodeString("dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877")
	hashes[8], _ = hex.DecodeString("b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59")
	hashes[9], _ = hex.DecodeString("95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c")
	hashes[10], _ = hex.DecodeString("2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908")

	var expected [][]byte = make([][]byte, 6)
	expected[0], _ = hex.DecodeString("8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd")
	expected[1], _ = hex.DecodeString("7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800")
	expected[2], _ = hex.DecodeString("ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7")
	expected[3], _ = hex.DecodeString("68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069")
	expected[4], _ = hex.DecodeString("43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27")
	expected[5], _ = hex.DecodeString("1796cd3ca4fef00236e07b723d3ed88e1ac433acaaa21da64c4b33c946cf3d10")

	results := MerkelParentLevel(hashes)

	if len(results) != len(expected) {
		t.Error()
	}

	for i := 0; i < len(results); i++ {

		if !bytes.Equal(expected[i], results[i]) {
			t.Error()
		}
	}
}

func TestMerkleRoot(t *testing.T) {

	var hashes [][]byte = make([][]byte, 12)
	hashes[0], _ = hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hashes[1], _ = hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	hashes[2], _ = hex.DecodeString("f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0")
	hashes[3], _ = hex.DecodeString("3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181")
	hashes[4], _ = hex.DecodeString("10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae")
	hashes[5], _ = hex.DecodeString("7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161")
	hashes[6], _ = hex.DecodeString("8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc")
	hashes[7], _ = hex.DecodeString("dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877")
	hashes[8], _ = hex.DecodeString("b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59")
	hashes[9], _ = hex.DecodeString("95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c")
	hashes[10], _ = hex.DecodeString("2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908")
	hashes[11], _ = hex.DecodeString("b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0")

	expected, _ := hex.DecodeString("acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6")

	result := MerkelRoot(hashes)

	if !bytes.Equal(expected, result) {
		t.Error()
	}
}
