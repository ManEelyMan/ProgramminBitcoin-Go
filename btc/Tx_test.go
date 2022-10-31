package btc_test

import (
	"bitcoin-go/btc"
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestTxParseVersion(t *testing.T) {
	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	buff := bytes.NewBuffer(b)

	tx := btc.ParseTx(buff, false)

	if tx.Version != 1 {
		t.Error()
	}
}

func TestTxParseInputs(t *testing.T) {
	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	buff := bytes.NewBuffer(b)
	tx := btc.ParseTx(buff, false)

	if len(tx.TxIns) != 1 {
		t.Error()
	}

	prevTx := utility.HexStringToBigInt("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	if tx.TxIns[0].PreviousTxHash.Cmp(prevTx) != 0 {
		t.Error()
	}

	if tx.TxIns[0].PreviousTxId != 0 {
		t.Error()
	}

	s, _ := hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	scriptData := bytes.NewBuffer(make([]byte, 0))
	tx.TxIns[0].ScriptSignature.Serialize(scriptData)

	if !bytes.Equal(scriptData.Bytes(), s) {
		t.Error()
	}

	if tx.TxIns[0].Sequence != 0xfffffffe {
		t.Error()
	}
}

func TestTxParseOutputs(t *testing.T) {

	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	reader := bytes.NewBuffer(b)
	tx := btc.ParseTx(reader, false)

	if len(tx.TxOuts) != 2 {
		t.Error()
	}

	if tx.TxOuts[0].Satoshis != 32454049 {
		t.Error()
	}

	writer := bytes.NewBuffer(make([]byte, 0))
	tx.TxOuts[0].ScriptPubKey.Serialize(writer)
	expected, _ := hex.DecodeString("1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac")

	if !bytes.Equal(writer.Bytes(), expected) {
		t.Error()
	}

	if tx.TxOuts[1].Satoshis != 10011545 {
		t.Error()
	}

	writer.Reset()
	tx.TxOuts[1].ScriptPubKey.Serialize(writer)
	expected, _ = hex.DecodeString("1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac")
	if !bytes.Equal(writer.Bytes(), expected) {
		t.Error()
	}
}

func TestTxParseLocktime(t *testing.T) {
	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	reader := bytes.NewBuffer(b)
	tx := btc.ParseTx(reader, false)

	if tx.LockTime != 410393 {
		t.Error()
	}
}

func TestTxFetcher(t *testing.T) {

	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	reader := bytes.NewBuffer(b)

	tx := btc.ParseTx(reader, false)
	f := btc.GetTxFetcher()

	txId1 := tx.TxIns[0].PreviousTxHash
	txIn1 := f.FetchById(txId1, false, false)

	// Make sure cache woks.
	txIn2 := f.FetchById(txId1, false, false)

	fmt.Print(txIn1, txIn2)
}

func TestTxCalcFee(t *testing.T) {

	b, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	reader := bytes.NewBuffer(b)
	tx := btc.ParseTx(reader, false)
	if tx.Fee() != 40000 {
		t.Error()
	}

	b, _ = hex.DecodeString("010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600")
	reader = bytes.NewBuffer(b)
	tx = btc.ParseTx(reader, false)
	if tx.Fee() != 140500 {
		t.Error()
	}
}

func TestTxSerialize(t *testing.T) {

	if !roundTripParseAndSerializationCheck("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600") {
		t.Error()
	}

	if !roundTripParseAndSerializationCheck("0100000002137c53f0fb48f83666fcfd2fe9f12d13e94ee109c5aeabbfa32bb9e02538f4cb000000006a47304402207e6009ad86367fc4b166bc80bf10cf1e78832a01e9bb491c6d126ee8aa436cb502200e29e6dd7708ed419cd5ba798981c960f0cc811b24e894bff072fea8074a7c4c012103bc9e7397f739c70f424aa7dcce9d2e521eb228b0ccba619cd6a0b9691da796a1ffffffff517472e77bc29ae59a914f55211f05024556812a2dd7d8df293265acd8330159010000006b483045022100f4bfdb0b3185c778cf28acbaf115376352f091ad9e27225e6f3f350b847579c702200d69177773cd2bb993a816a5ae08e77a6270cf46b33f8f79d45b0cd1244d9c4c0121031c0b0b95b522805ea9d0225b1946ecaeb1727c0b36c7e34165769fd8ed860bf5ffffffff027a958802000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88aca5515e00000000001976a914e82bd75c9c662c3f5700b33fec8a676b6e9391d588ac00000000") {
		t.Error()
	}

	if !roundTripParseAndSerializationCheck("0100000001c228021e1fee6f158cc506edea6bad7ffa421dd14fb7fd7e01c50cc9693e8dbe02000000fdfe0000483045022100c679944ff8f20373685e1122b581f64752c1d22c67f6f3ae26333aa9c3f43d730220793233401f87f640f9c39207349ffef42d0e27046755263c0a69c436ab07febc01483045022100eadc1c6e72f241c3e076a7109b8053db53987f3fcc99e3f88fc4e52dbfd5f3a202201f02cbff194c41e6f8da762e024a7ab85c1b1616b74720f13283043e9e99dab8014c69522102b0c7be446b92624112f3c7d4ffc214921c74c1cb891bf945c49fbe5981ee026b21039021c9391e328e0cb3b61ba05dcc5e122ab234e55d1502e59b10d8f588aea4632102f3bd8f64363066f35968bd82ed9c6e8afecbd6136311bb51e91204f614144e9b53aeffffffff05a08601000000000017a914081fbb6ec9d83104367eb1a6a59e2a92417d79298700350c00000000001976a914677345c7376dfda2c52ad9b6a153b643b6409a3788acc7f341160000000017a914234c15756b9599314c9299340eaabab7f1810d8287c02709000000000017a91469be3ca6195efcab5194e1530164ec47637d44308740420f00000000001976a91487fadba66b9e48c0c8082f33107fdb01970eb80388ac00000000") {
		t.Error()
	}

	if !roundTripParseAndSerializationCheck("0100000001b74780c0b9903472f84f8697a7449faebbfb1af659ecb8148ce8104347f3f72d010000006b483045022100bb8792c98141bcf4dab4fd4030743b4eff9edde59cec62380c60ffb90121ab7802204b439e3572b51382540c3b652b01327ee8b14cededc992fbc69b1e077a2c3f9f0121027c975c8bdc9717de310998494a2ae63f01b7a390bd34ef5b4c346fa717cba012ffffffff01a627c901000000001976a914af24b3f3e987c23528b366122a7ed2af199b36bc88ac00000000") {
		t.Error()
	}

	if !roundTripParseAndSerializationCheck("010000000367d54ded4c43569acbc213073fc63bfc49bf420391f0ab304758b16600a8ea88010000006a4730440220404b3bb28af45437c989328122aa6f4462021a0a2d4f20141ebe84e80edd72e202204184dd9d833d57246eaeed39021e9ab8c0546f3270bd9d2fc138a4bf161ea2310121039550662b907f788cc96708dc017aee0d407b74427f11e656b87f84146337f183feffffff5edf7dbc586b5fddace63a6614f5a731787c104d3c1c9225c4542db067d4296d010000006b483045022100b2335adb91e1ac3bb4e0479b54a9e7d4b765d9b646ca71e2547776c4e7e6bdfb02201fa8aaa4d2557768329befd61d4abda95668f88065df6eac6076e3e123c121eb012103b80229ec7a62793132ff432be0ecf21bca774ade18af7eaf2215febad0c4321ffeffffffdfa74eb50768daeb4beca2ca83d1732128d2439f9df9508efc8f7820718b4ae1000000006a47304402204818b29bed4a8ea4eb383f996389866a732b44d98f6342ecc25007ca472526fb0220496ed1213d63b7686f6936940e8f566f291bab211e6600c0f71e3659787b91fc0121036a30f9e6f645191c6216f84c21ae3b4f0aca0c4be987889276089cf9ef7a89d6feffffff028deb0f00000000001976a914cd0b3a22cd16e182291aa2708c41cb38de5a330788acc0e1e400000000001976a91424505f6d2f0fe7c4a3f4af32f50506034d89095d88ac43430600") {
		t.Error()
	}

	if !roundTripParseAndSerializationCheck("01000000012aa311f7789d362ceb2d802a98a703e0ac44815c021293633b80d08e67232e36010000006a4730440220142d8810ab29cac9199e6b570d47bd5ee402accf9d754cfa7de9b2e84e3997b402207a7d8c77c6a721bc64dba39eabe23e915c979683e621921c243bb35b3f538dfb01210371cb7d04e95471c4ea5c200e8c4729608754c74bee4e289bd66f431482407ec8feffffff02a08601000000000017a914fc7d096f19063ece361e2b309ec4da41fe4d789487f2798e00000000001976a914311b232c3400080eb2636edb8548b47f6835be7688ac31430600") {
		t.Error()
	}
}

func TestInputValue(t *testing.T) {
	txHash := utility.HexStringToBigInt("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	index := 0
	expected := 42505594
	txIn := btc.NewTxIn(txHash, uint32(index), nil, 0xffffffff)

	if txIn.Value(false) != uint64(expected) {
		t.Error()
	}
}

func TestInputPubKey(t *testing.T) {
	txHash := utility.HexStringToBigInt("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	index := 0
	txIn := btc.NewTxIn(txHash, uint32(index), nil, 0xffffffff)
	want, _ := hex.DecodeString("1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac")

	if bytes.Equal(want, txIn.ScriptPubKey(false).Data) {
		t.Error()
	}
}

func TestFee(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	reader := bytes.NewBuffer(rawTx)
	tx := btc.ParseTx(reader, false)

	if tx.Fee() != 40000 {
		t.Error()
	}

	rawTx, _ = hex.DecodeString("010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600")
	reader = bytes.NewBuffer(rawTx)
	tx = btc.ParseTx(reader, false)

	if tx.Fee() != 140500 {
		t.Error()
	}
}

func TestSigHash(t *testing.T) {
	fetch := btc.GetTxFetcher()
	tx := fetch.FetchById(utility.HexStringToBigInt("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"), false, false)
	expected, _ := hex.DecodeString("27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6")
	if bytes.Equal(tx.SigHash(0), expected) {
		t.Error()
	}
}

func TestVerifyP2PKH(t *testing.T) {
	fetcher := btc.GetTxFetcher()
	tx1 := fetcher.FetchById(utility.HexStringToBigInt("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"), false, false)
	if !tx1.Verify() {
		t.Error()
	}

	tx2 := fetcher.FetchById(utility.HexStringToBigInt("5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2"), true, false)
	if !tx2.Verify() {
		t.Error()
	}

}

func roundTripParseAndSerializationCheck(encodedHex string) bool {
	b, _ := hex.DecodeString(encodedHex)
	reader := bytes.NewBuffer(b)
	tx := btc.ParseTx(reader, false)

	writer := bytes.NewBuffer(make([]byte, 0))
	tx.Serialize(writer, -1)

	b2 := writer.Bytes()

	return bytes.Equal(b, b2)
}
