package core

import (
	"errors"
	"fmt"

	ctypes "github.com/creatachain/augusteum/rpc/core/types"
	rpctypes "github.com/creatachain/augusteum/rpc/jsonrpc/types"
	"github.com/creatachain/augusteum/types"
)

// BroadcastEvidence broadcasts evidence of the misbehavior.
// More: https://docs.augusteum.com/master/rpc/#/Info/broadcast_evidence
func BroadcastEvidence(ctx *rpctypes.Context, ev types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	if ev == nil {
		return nil, errors.New("no evidence was provided")
	}

	if err := ev.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("evidence.ValidateBasic failed: %w", err)
	}

	if err := env.EvidencePool.AddEvidence(ev); err != nil {
		return nil, fmt.Errorf("failed to add evidence: %w", err)
	}
	return &ctypes.ResultBroadcastEvidence{Hash: ev.Hash()}, nil
}
