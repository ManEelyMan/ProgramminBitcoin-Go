package btc

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

const SIGHASH_ALL = 1

type Tx struct {
	Version  uint32
	TxIns    []TxIn
	TxOuts   []TxOut
	LockTime uint32
	TestNet  bool
}

func NewTx(version uint32, txIns []TxIn, txOuts []TxOut, lockTime uint32, testNet bool) Tx {
	return Tx{Version: version, TxIns: txIns, TxOuts: txOuts, LockTime: lockTime, TestNet: testNet}
}

func ParseTx(reader io.Reader, testnet bool) Tx {
	version := utility.ReadUint32(reader, true)

	txInCount := utility.ReadVarInt(reader)

	txIns := make([]TxIn, txInCount)
	for i := (uint64)(0); i < txInCount; i++ {
		txIns[i] = ParseTxIn(reader)
	}

	txOutCount := utility.ReadVarInt(reader)

	txOuts := make([]TxOut, txOutCount)
	for i := (uint64)(0); i < txOutCount; i++ {
		txOuts[i] = ParseTxOut(reader)
	}

	lockTime := utility.ReadUint32(reader, true)

	return Tx{Version: version, TxIns: txIns, TxOuts: txOuts, LockTime: lockTime, TestNet: testnet}
}

func (tx *Tx) Id() string {
	return hex.EncodeToString(tx.Hash())
}

func (tx *Tx) Hash() []byte {
	h := tx.SigHash(-1)
	return utility.ReverseBytes(h)
}

func (tx *Tx) SigHash(txSigHash int) []byte {

	buff := bytes.NewBuffer(make([]byte, 0))
	tx.Serialize(buff, txSigHash)
	if txSigHash >= 0 {
		utility.WriteUint32(buff, SIGHASH_ALL, true)
	}
	return utility.Hash256(buff.Bytes())
}

func (tx *Tx) Serialize(writer io.Writer, txSigHash int) {

	utility.WriteUint32(writer, tx.Version, true)
	utility.WriteVarInt(writer, (uint64)(len(tx.TxIns)))

	for i, txin := range tx.TxIns {
		txin.Serialize(writer, i == txSigHash, tx.TestNet)
	}

	utility.WriteVarInt(writer, (uint64)(len(tx.TxOuts)))

	for _, txout := range tx.TxOuts {
		txout.Serialize(writer)
	}

	utility.WriteUint32(writer, tx.LockTime, true)
}

func (tx *Tx) Fee() int64 {
	var input_sum uint64 = 0
	var output_sum uint64 = 0

	for _, in := range tx.TxIns {
		input_sum += in.Value(tx.TestNet)
	}

	for _, out := range tx.TxOuts {
		output_sum += out.Satoshis
	}

	return int64(input_sum - output_sum)
}

func (tx *Tx) Verify() bool {
	if tx.Fee() < 0 {
		return false
	}

	for i := range tx.TxIns {
		if !tx.VerifyInput(i) {
			return false
		}
	}

	return true
}

func (tx *Tx) VerifyInput(index int) bool {
	txIn := tx.TxIns[index]
	scriptPubKey := txIn.ScriptPubKey(tx.TestNet)
	z := tx.SigHash(index)

	hash := big.NewInt(0)
	hash.SetBytes(z)
	fmt.Printf("Hash: %+v\n", hash)
	exec := NewScriptExecutor(&scriptPubKey, txIn.ScriptSignature, hash)
	return exec.Execute()
}
