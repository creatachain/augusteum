package kv

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"testing"

	dbm "github.com/creatachain/tm-db"

	"github.com/creatachain/augusteum/libs/pubsub/query"
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/types"
)

func BenchmarkTxSearch(b *testing.B) {
	dbDir, err := ioutil.TempDir("", "benchmark_tx_search_test")
	if err != nil {
		b.Errorf("failed to create temporary directory: %s", err)
	}

	db, err := dbm.NewGoLevelDB("benchmark_tx_search_test", dbDir)
	if err != nil {
		b.Errorf("failed to create database: %s", err)
	}

	indexer := NewTxIndex(db)

	for i := 0; i < 35000; i++ {
		events := []msm.Event{
			{
				Type: "transfer",
				Attributes: []msm.EventAttribute{
					{Key: []byte("address"), Value: []byte(fmt.Sprintf("address_%d", i%100)), Index: true},
					{Key: []byte("amount"), Value: []byte("50"), Index: true},
				},
			},
		}

		txBz := make([]byte, 8)
		if _, err := rand.Read(txBz); err != nil {
			b.Errorf("failed produce random bytes: %s", err)
		}

		txResult := &msm.TxResult{
			Height: int64(i),
			Index:  0,
			Tx:     types.Tx(string(txBz)),
			Result: msm.ResponseDeliverTx{
				Data:   []byte{0},
				Code:   msm.CodeTypeOK,
				Log:    "",
				Events: events,
			},
		}

		if err := indexer.Index(txResult); err != nil {
			b.Errorf("failed to index tx: %s", err)
		}
	}

	txQuery := query.MustParse("transfer.address = 'address_43' AND transfer.amount = 50")

	b.ResetTimer()

	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		if _, err := indexer.Search(ctx, txQuery); err != nil {
			b.Errorf("failed to query for txs: %s", err)
		}
	}
}
