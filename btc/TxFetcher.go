package btc

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type TxFetcher struct {
	cache map[[32]byte]Tx
}

var once sync.Once
var singleton TxFetcher

func GetTxFetcher() TxFetcher {

	once.Do(func() {
		singleton = TxFetcher{}
		singleton.cache = make(map[[32]byte]Tx)
	})

	return singleton
}

func (f *TxFetcher) FetchById(txId [32]byte, testNet bool, fresh bool) Tx {
	tx, ok := singleton.cache[txId]
	if !ok || fresh {
		tx = f.fetchTransaction(txId, testNet)
		f.cache[txId] = tx // TODO: Create a disk-persisting cache.
	}

	return tx
}

func (f *TxFetcher) fetchTransaction(txId [32]byte, testNet bool) Tx {
	url := utility.IIF(testNet, "https://blockstream.info/testnet/api/tx/%v/hex", "https://blockstream.info/api/tx/%v/hex").(string)

	url = fmt.Sprintf(url, hex.EncodeToString(txId[:]))

	client := http.DefaultClient
	rsp, err := client.Get(url)

	if err != nil {
		panic(err)
	}

	// The body comes back as an ASCII string of the hex.
	b, _ := io.ReadAll(rsp.Body)
	s := (string)(b)
	b, _ = hex.DecodeString(s)

	/*
		TODO: Figure out why they do this:
		if raw[4] == 0:
			raw = raw[:4] + raw[6:]
			tx = Tx.parse(BytesIO(raw), testnet=testnet)
			tx.locktime = little_endian_to_int(raw[-4:])

	*/

	reader := bytes.NewBuffer(b)
	tx := ParseTx(reader, testNet)
	txHash := tx.Hash()

	if !bytes.Equal(txHash, txId[:]) {
		fmt.Println("Sanity check broke: Id mismatch between requested and received Transaction!")
	}
	return tx
}
