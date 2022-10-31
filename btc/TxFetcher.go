package btc

import (
	"bitcoin-go/utility"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
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

func (f *TxFetcher) FetchById(txId *big.Int, testNet bool, fresh bool) Tx {
	var bytes [32]byte
	copy(bytes[0:], txId.Bytes())

	tx, ok := singleton.cache[bytes]
	if !ok || fresh {
		tx = f.fetchTransaction(txId, testNet)
		f.cache[bytes] = tx // TODO: Create a disk-persisting cache.
	}

	return tx
}

func (f *TxFetcher) fetchTransaction(txId *big.Int, testNet bool) Tx {
	url := utility.IIF(testNet, "https://blockstream.info/testnet/api/tx/%v/hex", "https://blockstream.info/api/tx/%v/hex").(string)

	idBytes := txId.Bytes()
	url = fmt.Sprintf(url, hex.EncodeToString(idBytes))

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

	if !bytes.Equal(tx.Hash(), idBytes) {
		fmt.Println("Sanity check broke: Id mismatch between requested and received Transaction!")
	}
	return tx
}
