package state_test

import (
	"bytes"
	"fmt"
	"time"

	dbm "github.com/creatachain/tm-db"

	"github.com/creatachain/augusteum/crypto"
	"github.com/creatachain/augusteum/crypto/ed25519"
	tmrand "github.com/creatachain/augusteum/libs/rand"
	msm "github.com/creatachain/augusteum/msm/types"
	tmstate "github.com/creatachain/augusteum/proto/augusteum/state"
	tmproto "github.com/creatachain/augusteum/proto/augusteum/types"
	"github.com/creatachain/augusteum/proxy"
	sm "github.com/creatachain/augusteum/state"
	"github.com/creatachain/augusteum/types"
	tmtime "github.com/creatachain/augusteum/types/time"
)

type paramsChangeTestCase struct {
	height int64
	params tmproto.ConsensusParams
}

func newTestApp() proxy.AppConns {
	app := &testApp{}
	cc := proxy.NewLocalClientCreator(app)
	return proxy.NewAppConns(cc)
}

func makeAndCommitGoodBlock(
	state sm.State,
	height int64,
	lastCommit *types.Commit,
	proposerAddr []byte,
	blockExec *sm.BlockExecutor,
	privVals map[string]types.PrivValidator,
	evidence []types.Evidence) (sm.State, types.BlockID, *types.Commit, error) {
	// A good block passes
	state, blockID, err := makeAndApplyGoodBlock(state, height, lastCommit, proposerAddr, blockExec, evidence)
	if err != nil {
		return state, types.BlockID{}, nil, err
	}

	// Simulate a lastCommit for this block from all validators for the next height
	commit, err := makeValidCommit(height, blockID, state.Validators, privVals)
	if err != nil {
		return state, types.BlockID{}, nil, err
	}
	return state, blockID, commit, nil
}

func makeAndApplyGoodBlock(state sm.State, height int64, lastCommit *types.Commit, proposerAddr []byte,
	blockExec *sm.BlockExecutor, evidence []types.Evidence) (sm.State, types.BlockID, error) {
	block, _ := state.MakeBlock(height, makeTxs(height), lastCommit, evidence, proposerAddr)
	if err := blockExec.ValidateBlock(state, block); err != nil {
		return state, types.BlockID{}, err
	}
	blockID := types.BlockID{Hash: block.Hash(),
		PartSetHeader: types.PartSetHeader{Total: 3, Hash: tmrand.Bytes(32)}}
	state, _, err := blockExec.ApplyBlock(state, blockID, block)
	if err != nil {
		return state, types.BlockID{}, err
	}
	return state, blockID, nil
}

func makeValidCommit(
	height int64,
	blockID types.BlockID,
	vals *types.ValidatorSet,
	privVals map[string]types.PrivValidator,
) (*types.Commit, error) {
	sigs := make([]types.CommitSig, 0)
	for i := 0; i < vals.Size(); i++ {
		_, val := vals.GetByIndex(int32(i))
		vote, err := types.MakeVote(height, blockID, vals, privVals[val.Address.String()], chainID, time.Now())
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, vote.CommitSig())
	}
	return types.NewCommit(height, 0, blockID, sigs), nil
}

// make some bogus txs
func makeTxs(height int64) (txs []types.Tx) {
	for i := 0; i < nTxsPerBlock; i++ {
		txs = append(txs, types.Tx([]byte{byte(height), byte(i)}))
	}
	return txs
}

func makeState(nVals, height int) (sm.State, dbm.DB, map[string]types.PrivValidator) {
	vals := make([]types.GenesisValidator, nVals)
	privVals := make(map[string]types.PrivValidator, nVals)
	for i := 0; i < nVals; i++ {
		secret := []byte(fmt.Sprintf("test%d", i))
		pk := ed25519.GenPrivKeyFromSecret(secret)
		valAddr := pk.PubKey().Address()
		vals[i] = types.GenesisValidator{
			Address: valAddr,
			PubKey:  pk.PubKey(),
			Power:   1000,
			Name:    fmt.Sprintf("test%d", i),
		}
		privVals[valAddr.String()] = types.NewMockPVWithParams(pk, false, false)
	}
	s, _ := sm.MakeGenesisState(&types.GenesisDoc{
		ChainID:    chainID,
		Validators: vals,
		AppHash:    nil,
	})

	stateDB := dbm.NewMemDB()
	stateStore := sm.NewStore(stateDB)
	if err := stateStore.Save(s); err != nil {
		panic(err)
	}

	for i := 1; i < height; i++ {
		s.LastBlockHeight++
		s.LastValidators = s.Validators.Copy()
		if err := stateStore.Save(s); err != nil {
			panic(err)
		}
	}

	return s, stateDB, privVals
}

func makeBlock(state sm.State, height int64) *types.Block {
	block, _ := state.MakeBlock(
		height,
		makeTxs(state.LastBlockHeight),
		new(types.Commit),
		nil,
		state.Validators.GetProposer().Address,
	)
	return block
}

func genValSet(size int) *types.ValidatorSet {
	vals := make([]*types.Validator, size)
	for i := 0; i < size; i++ {
		vals[i] = types.NewValidator(ed25519.GenPrivKey().PubKey(), 10)
	}
	return types.NewValidatorSet(vals)
}

func makeHeaderPartsResponsesValPubKeyChange(
	state sm.State,
	pubkey crypto.PubKey,
) (types.Header, types.BlockID, *tmstate.MSMResponses) {

	block := makeBlock(state, state.LastBlockHeight+1)
	msmResponses := &tmstate.MSMResponses{
		BeginBlock: &msm.ResponseBeginBlock{},
		EndBlock:   &msm.ResponseEndBlock{ValidatorUpdates: nil},
	}
	// If the pubkey is new, remove the old and add the new.
	_, val := state.NextValidators.GetByIndex(0)
	if !bytes.Equal(pubkey.Bytes(), val.PubKey.Bytes()) {
		msmResponses.EndBlock = &msm.ResponseEndBlock{
			ValidatorUpdates: []msm.ValidatorUpdate{
				types.TM2PB.NewValidatorUpdate(val.PubKey, 0),
				types.TM2PB.NewValidatorUpdate(pubkey, 10),
			},
		}
	}

	return block.Header, types.BlockID{Hash: block.Hash(), PartSetHeader: types.PartSetHeader{}}, msmResponses
}

func makeHeaderPartsResponsesValPowerChange(
	state sm.State,
	power int64,
) (types.Header, types.BlockID, *tmstate.MSMResponses) {

	block := makeBlock(state, state.LastBlockHeight+1)
	msmResponses := &tmstate.MSMResponses{
		BeginBlock: &msm.ResponseBeginBlock{},
		EndBlock:   &msm.ResponseEndBlock{ValidatorUpdates: nil},
	}

	// If the pubkey is new, remove the old and add the new.
	_, val := state.NextValidators.GetByIndex(0)
	if val.VotingPower != power {
		msmResponses.EndBlock = &msm.ResponseEndBlock{
			ValidatorUpdates: []msm.ValidatorUpdate{
				types.TM2PB.NewValidatorUpdate(val.PubKey, power),
			},
		}
	}

	return block.Header, types.BlockID{Hash: block.Hash(), PartSetHeader: types.PartSetHeader{}}, msmResponses
}

func makeHeaderPartsResponsesParams(
	state sm.State,
	params tmproto.ConsensusParams,
) (types.Header, types.BlockID, *tmstate.MSMResponses) {

	block := makeBlock(state, state.LastBlockHeight+1)
	msmResponses := &tmstate.MSMResponses{
		BeginBlock: &msm.ResponseBeginBlock{},
		EndBlock:   &msm.ResponseEndBlock{ConsensusParamUpdates: types.TM2PB.ConsensusParams(&params)},
	}
	return block.Header, types.BlockID{Hash: block.Hash(), PartSetHeader: types.PartSetHeader{}}, msmResponses
}

func randomGenesisDoc() *types.GenesisDoc {
	pubkey := ed25519.GenPrivKey().PubKey()
	return &types.GenesisDoc{
		GenesisTime: tmtime.Now(),
		ChainID:     "abc",
		Validators: []types.GenesisValidator{
			{
				Address: pubkey.Address(),
				PubKey:  pubkey,
				Power:   10,
				Name:    "myval",
			},
		},
		ConsensusParams: types.DefaultConsensusParams(),
	}
}

//----------------------------------------------------------------------------

type testApp struct {
	msm.BaseApplication

	CommitVotes         []msm.VoteInfo
	ByzantineValidators []msm.Evidence
	ValidatorUpdates    []msm.ValidatorUpdate
}

var _ msm.Application = (*testApp)(nil)

func (app *testApp) Info(req msm.RequestInfo) (resInfo msm.ResponseInfo) {
	return msm.ResponseInfo{}
}

func (app *testApp) BeginBlock(req msm.RequestBeginBlock) msm.ResponseBeginBlock {
	app.CommitVotes = req.LastCommitInfo.Votes
	app.ByzantineValidators = req.ByzantineValidators
	return msm.ResponseBeginBlock{}
}

func (app *testApp) EndBlock(req msm.RequestEndBlock) msm.ResponseEndBlock {
	return msm.ResponseEndBlock{
		ValidatorUpdates: app.ValidatorUpdates,
		ConsensusParamUpdates: &msm.ConsensusParams{
			Version: &tmproto.VersionParams{
				AppVersion: 1}}}
}

func (app *testApp) DeliverTx(req msm.RequestDeliverTx) msm.ResponseDeliverTx {
	return msm.ResponseDeliverTx{Events: []msm.Event{}}
}

func (app *testApp) CheckTx(req msm.RequestCheckTx) msm.ResponseCheckTx {
	return msm.ResponseCheckTx{}
}

func (app *testApp) Commit() msm.ResponseCommit {
	return msm.ResponseCommit{RetainHeight: 1}
}

func (app *testApp) Query(reqQuery msm.RequestQuery) (resQuery msm.ResponseQuery) {
	return
}
