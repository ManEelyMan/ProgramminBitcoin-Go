package btc

import (
	"bitcoin-go/utility"
	"io"
	"math/big"
)

type TxIn struct {
	PreviousTxHash  *big.Int
	PreviousTxId    uint32
	ScriptSignature *Script
	Sequence        uint32
}

func NewTxIn(prevTxHash *big.Int, prevTxId uint32, scriptSig *Script, sequence uint32) TxIn {
	return TxIn{PreviousTxHash: prevTxHash, PreviousTxId: prevTxId, ScriptSignature: scriptSig, Sequence: sequence}
}

func ParseTxIn(reader io.Reader) TxIn {
	prevTx := utility.ReadBigInt(reader, true)
	prevTxId := utility.ReadUint32(reader, true)
	script := ParseScript(reader)
	sequence := utility.ReadUint32(reader, true)
	return TxIn{PreviousTxHash: prevTx, PreviousTxId: prevTxId, ScriptSignature: &script, Sequence: sequence}
}

func (txin *TxIn) Serialize(writer io.Writer, sigHash bool, testNet bool, redeemScript *Script) {
	utility.WriteBigInt(writer, txin.PreviousTxHash, true)
	utility.WriteUint32(writer, txin.PreviousTxId, true)

	if sigHash {
		spk := txin.ScriptPubKey(testNet)
		clone := Script{}

		if redeemScript != nil {
			clone.RawData = redeemScript.RawData
		} else {
			clone.RawData = append(clone.RawData, spk.RawData...)

		}
		clone.Serialize(writer)
	} else {
		txin.ScriptSignature.Serialize(writer)
	}

	utility.WriteUint32(writer, txin.Sequence, true)
}

func (txin *TxIn) Value(testNet bool) uint64 {
	tx := txin.PreviousTx(testNet)
	return tx.TxOuts[txin.PreviousTxId].Satoshis
}

func (txin *TxIn) PreviousTx(testNet bool) Tx {
	f := GetTxFetcher()
	tx := f.FetchById(txin.PreviousTxHash, testNet, false)
	return tx
}

func (txin *TxIn) ScriptPubKey(testNet bool) Script {
	return txin.PreviousTx(testNet).TxOuts[txin.PreviousTxId].ScriptPubKey
}
