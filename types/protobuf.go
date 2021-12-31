package types

import (
	"github.com/creatachain/augusteum/crypto"
	"github.com/creatachain/augusteum/crypto/ed25519"
	cryptoenc "github.com/creatachain/augusteum/crypto/encoding"
	"github.com/creatachain/augusteum/crypto/secp256k1"
	msm "github.com/creatachain/augusteum/msm/types"
	tmproto "github.com/creatachain/augusteum/proto/augusteum/types"
)

//-------------------------------------------------------
// Use strings to distinguish types in MSM messages

const (
	MSMPubKeyTypeEd25519   = ed25519.KeyType
	MSMPubKeyTypeSecp256k1 = secp256k1.KeyType
)

// TODO: Make non-global by allowing for registration of more pubkey types

var MSMPubKeyTypesToNames = map[string]string{
	MSMPubKeyTypeEd25519:   ed25519.PubKeyName,
	MSMPubKeyTypeSecp256k1: secp256k1.PubKeyName,
}

//-------------------------------------------------------

// TM2PB is used for converting Augusteum MSM to protobuf MSM.
// UNSTABLE
var TM2PB = tm2pb{}

type tm2pb struct{}

func (tm2pb) Header(header *Header) tmproto.Header {
	return tmproto.Header{
		Version: header.Version,
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,

		LastBlockId: header.LastBlockID.ToProto(),

		LastCommitHash: header.LastCommitHash,
		DataHash:       header.DataHash,

		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,

		EvidenceHash:    header.EvidenceHash,
		ProposerAddress: header.ProposerAddress,
	}
}

func (tm2pb) Validator(val *Validator) msm.Validator {
	return msm.Validator{
		Address: val.PubKey.Address(),
		Power:   val.VotingPower,
	}
}

func (tm2pb) BlockID(blockID BlockID) tmproto.BlockID {
	return tmproto.BlockID{
		Hash:          blockID.Hash,
		PartSetHeader: TM2PB.PartSetHeader(blockID.PartSetHeader),
	}
}

func (tm2pb) PartSetHeader(header PartSetHeader) tmproto.PartSetHeader {
	return tmproto.PartSetHeader{
		Total: header.Total,
		Hash:  header.Hash,
	}
}

// XXX: panics on unknown pubkey type
func (tm2pb) ValidatorUpdate(val *Validator) msm.ValidatorUpdate {
	pk, err := cryptoenc.PubKeyToProto(val.PubKey)
	if err != nil {
		panic(err)
	}
	return msm.ValidatorUpdate{
		PubKey: pk,
		Power:  val.VotingPower,
	}
}

// XXX: panics on nil or unknown pubkey type
func (tm2pb) ValidatorUpdates(vals *ValidatorSet) []msm.ValidatorUpdate {
	validators := make([]msm.ValidatorUpdate, vals.Size())
	for i, val := range vals.Validators {
		validators[i] = TM2PB.ValidatorUpdate(val)
	}
	return validators
}

func (tm2pb) ConsensusParams(params *tmproto.ConsensusParams) *msm.ConsensusParams {
	return &msm.ConsensusParams{
		Block: &msm.BlockParams{
			MaxBytes: params.Block.MaxBytes,
			MaxGas:   params.Block.MaxGas,
		},
		Evidence:  &params.Evidence,
		Validator: &params.Validator,
	}
}

// XXX: panics on nil or unknown pubkey type
func (tm2pb) NewValidatorUpdate(pubkey crypto.PubKey, power int64) msm.ValidatorUpdate {
	pubkeyMSM, err := cryptoenc.PubKeyToProto(pubkey)
	if err != nil {
		panic(err)
	}
	return msm.ValidatorUpdate{
		PubKey: pubkeyMSM,
		Power:  power,
	}
}

//----------------------------------------------------------------------------

// PB2TM is used for converting protobuf MSM to Augusteum MSM.
// UNSTABLE
var PB2TM = pb2tm{}

type pb2tm struct{}

func (pb2tm) ValidatorUpdates(vals []msm.ValidatorUpdate) ([]*Validator, error) {
	tmVals := make([]*Validator, len(vals))
	for i, v := range vals {
		pub, err := cryptoenc.PubKeyFromProto(v.PubKey)
		if err != nil {
			return nil, err
		}
		tmVals[i] = NewValidator(pub, v.Power)
	}
	return tmVals, nil
}
