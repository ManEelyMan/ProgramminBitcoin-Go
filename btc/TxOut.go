package btc

import (
	"bitcoin-go/utility"
	"io"
)

type TxOut struct {
	Satoshis     uint64
	ScriptPubKey Script
}

func NewTxOut(satoshis uint64, scriptPubKey Script) TxOut {
	return TxOut{Satoshis: satoshis, ScriptPubKey: scriptPubKey}
}

func ParseTxOut(reader io.Reader) TxOut {
	sats := utility.ReadUint64(reader, true)
	scriptPubKey := ParseScript(reader)
	return TxOut{Satoshis: sats, ScriptPubKey: scriptPubKey}
}

func (txout *TxOut) Serialize(writer io.Writer) {
	utility.WriteUint64(writer, txout.Satoshis, true)
	txout.ScriptPubKey.Serialize(writer)
}
