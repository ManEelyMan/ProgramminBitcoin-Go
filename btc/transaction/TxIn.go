package transaction

import (
	"bitcoin-go/utility"
	"io"
)

type TxIn struct {
	PreviousTxHash  [32]byte
	PreviousTxId    uint32
	ScriptSignature *Script
	Sequence        uint32
}

func NewTxIn(prevTxHash [32]byte, prevTxId uint32, scriptSig *Script, sequence uint32) TxIn {
	return TxIn{PreviousTxHash: prevTxHash, PreviousTxId: prevTxId, ScriptSignature: scriptSig, Sequence: sequence}
}

func ParseTxIn(reader io.Reader) TxIn {
	prevTx, _ := utility.ReadBytes(reader, 32)
	prevTx = utility.ReverseBytes(prevTx)
	prevTxId := utility.ReadUint32(reader, true)
	script := ParseScript(reader)
	sequence := utility.ReadUint32(reader, true)
	txIn := TxIn{PreviousTxId: prevTxId, ScriptSignature: &script, Sequence: sequence}
	copy(txIn.PreviousTxHash[:], prevTx)
	return txIn
}

func (txin *TxIn) Serialize(writer io.Writer, sigHash bool, testNet bool, redeemScript *Script) {

	var reversed [32]byte
	copy(reversed[:], txin.PreviousTxHash[:])
	writer.Write(utility.ReverseBytes(reversed[:]))
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
	prevTx := txin.PreviousTx(testNet)
	prevTxOutput := prevTx.TxOuts[txin.PreviousTxId]
	return prevTxOutput.ScriptPubKey
}
