package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/creatachain/augusteum/crypto"
	"github.com/creatachain/augusteum/crypto/ed25519"
	cryptoenc "github.com/creatachain/augusteum/crypto/encoding"
	msm "github.com/creatachain/augusteum/msm/types"
)

func TestMSMPubKey(t *testing.T) {
	pkEd := ed25519.GenPrivKey().PubKey()
	err := testMSMPubKey(t, pkEd, MSMPubKeyTypeEd25519)
	assert.NoError(t, err)
}

func testMSMPubKey(t *testing.T, pk crypto.PubKey, typeStr string) error {
	msmPubKey, err := cryptoenc.PubKeyToProto(pk)
	require.NoError(t, err)
	pk2, err := cryptoenc.PubKeyFromProto(msmPubKey)
	require.NoError(t, err)
	require.Equal(t, pk, pk2)
	return nil
}

func TestMSMValidators(t *testing.T) {
	pkEd := ed25519.GenPrivKey().PubKey()

	// correct validator
	tmValExpected := NewValidator(pkEd, 10)

	tmVal := NewValidator(pkEd, 10)

	msmVal := TM2PB.ValidatorUpdate(tmVal)
	tmVals, err := PB2TM.ValidatorUpdates([]msm.ValidatorUpdate{msmVal})
	assert.Nil(t, err)
	assert.Equal(t, tmValExpected, tmVals[0])

	msmVals := TM2PB.ValidatorUpdates(NewValidatorSet(tmVals))
	assert.Equal(t, []msm.ValidatorUpdate{msmVal}, msmVals)

	// val with address
	tmVal.Address = pkEd.Address()

	msmVal = TM2PB.ValidatorUpdate(tmVal)
	tmVals, err = PB2TM.ValidatorUpdates([]msm.ValidatorUpdate{msmVal})
	assert.Nil(t, err)
	assert.Equal(t, tmValExpected, tmVals[0])
}

func TestMSMConsensusParams(t *testing.T) {
	cp := DefaultConsensusParams()
	msmCP := TM2PB.ConsensusParams(cp)
	cp2 := UpdateConsensusParams(*cp, msmCP)

	assert.Equal(t, *cp, cp2)
}

type pubKeyEddie struct{}

func (pubKeyEddie) Address() Address                            { return []byte{} }
func (pubKeyEddie) Bytes() []byte                               { return []byte{} }
func (pubKeyEddie) VerifySignature(msg []byte, sig []byte) bool { return false }
func (pubKeyEddie) Equals(crypto.PubKey) bool                   { return false }
func (pubKeyEddie) String() string                              { return "" }
func (pubKeyEddie) Type() string                                { return "pubKeyEddie" }

func TestMSMValidatorFromPubKeyAndPower(t *testing.T) {
	pubkey := ed25519.GenPrivKey().PubKey()

	msmVal := TM2PB.NewValidatorUpdate(pubkey, 10)
	assert.Equal(t, int64(10), msmVal.Power)

	assert.Panics(t, func() { TM2PB.NewValidatorUpdate(nil, 10) })
	assert.Panics(t, func() { TM2PB.NewValidatorUpdate(pubKeyEddie{}, 10) })
}

func TestMSMValidatorWithoutPubKey(t *testing.T) {
	pkEd := ed25519.GenPrivKey().PubKey()

	msmVal := TM2PB.Validator(NewValidator(pkEd, 10))

	// pubkey must be nil
	tmValExpected := msm.Validator{
		Address: pkEd.Address(),
		Power:   10,
	}

	assert.Equal(t, tmValExpected, msmVal)
}
