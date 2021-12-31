package types

import (
	"github.com/creatachain/augusteum/crypto/merkle"
	msm "github.com/creatachain/augusteum/msm/types"
)

// MSMResults wraps the deliver tx results to return a proof.
type MSMResults []*msm.ResponseDeliverTx

// NewResults strips non-deterministic fields from ResponseDeliverTx responses
// and returns MSMResults.
func NewResults(responses []*msm.ResponseDeliverTx) MSMResults {
	res := make(MSMResults, len(responses))
	for i, d := range responses {
		res[i] = deterministicResponseDeliverTx(d)
	}
	return res
}

// Hash returns a merkle hash of all results.
func (a MSMResults) Hash() []byte {
	return merkle.HashFromByteSlices(a.toByteSlices())
}

// ProveResult returns a merkle proof of one result from the set
func (a MSMResults) ProveResult(i int) merkle.Proof {
	_, proofs := merkle.ProofsFromByteSlices(a.toByteSlices())
	return *proofs[i]
}

func (a MSMResults) toByteSlices() [][]byte {
	l := len(a)
	bzs := make([][]byte, l)
	for i := 0; i < l; i++ {
		bz, err := a[i].Marshal()
		if err != nil {
			panic(err)
		}
		bzs[i] = bz
	}
	return bzs
}

// deterministicResponseDeliverTx strips non-deterministic fields from
// ResponseDeliverTx and returns another ResponseDeliverTx.
func deterministicResponseDeliverTx(response *msm.ResponseDeliverTx) *msm.ResponseDeliverTx {
	return &msm.ResponseDeliverTx{
		Code:      response.Code,
		Data:      response.Data,
		GasWanted: response.GasWanted,
		GasUsed:   response.GasUsed,
	}
}
