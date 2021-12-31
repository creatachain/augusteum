package null

import (
	"context"
	"errors"

	"github.com/creatachain/augusteum/libs/pubsub/query"
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/state/txindex"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

// TxIndex acts as a /dev/null.
type TxIndex struct{}

// Get on a TxIndex is disabled and panics when invoked.
func (txi *TxIndex) Get(hash []byte) (*msm.TxResult, error) {
	return nil, errors.New(`indexing is disabled (set 'tx_index = "kv"' in config)`)
}

// AddBatch is a noop and always returns nil.
func (txi *TxIndex) AddBatch(batch *txindex.Batch) error {
	return nil
}

// Index is a noop and always returns nil.
func (txi *TxIndex) Index(result *msm.TxResult) error {
	return nil
}

func (txi *TxIndex) Search(ctx context.Context, q *query.Query) ([]*msm.TxResult, error) {
	return []*msm.TxResult{}, nil
}
