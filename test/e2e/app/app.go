package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/creatachain/augusteum/libs/log"
	"github.com/creatachain/augusteum/msm/example/code"
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/version"
)

// Application is an MSM application for use by end-to-end tests. It is a
// simple key/value store for strings, storing data in memory and persisting
// to disk as JSON, taking state sync snapshots if requested.
type Application struct {
	msm.BaseApplication
	logger          log.Logger
	state           *State
	snapshots       *SnapshotStore
	cfg             *Config
	restoreSnapshot *msm.Snapshot
	restoreChunks   [][]byte
}

// NewApplication creates the application.
func NewApplication(cfg *Config) (*Application, error) {
	state, err := NewState(filepath.Join(cfg.Dir, "state.json"), cfg.PersistInterval)
	if err != nil {
		return nil, err
	}
	snapshots, err := NewSnapshotStore(filepath.Join(cfg.Dir, "snapshots"))
	if err != nil {
		return nil, err
	}
	return &Application{
		logger:    log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		state:     state,
		snapshots: snapshots,
		cfg:       cfg,
	}, nil
}

// Info implements MSM.
func (app *Application) Info(req msm.RequestInfo) msm.ResponseInfo {
	return msm.ResponseInfo{
		Version:          version.MSMVersion,
		AppVersion:       1,
		LastBlockHeight:  int64(app.state.Height),
		LastBlockAppHash: app.state.Hash,
	}
}

// Info implements MSM.
func (app *Application) InitChain(req msm.RequestInitChain) msm.ResponseInitChain {
	var err error
	app.state.initialHeight = uint64(req.InitialHeight)
	if len(req.AppStateBytes) > 0 {
		err = app.state.Import(0, req.AppStateBytes)
		if err != nil {
			panic(err)
		}
	}
	resp := msm.ResponseInitChain{
		AppHash: app.state.Hash,
	}
	if resp.Validators, err = app.validatorUpdates(0); err != nil {
		panic(err)
	}
	return resp
}

// CheckTx implements MSM.
func (app *Application) CheckTx(req msm.RequestCheckTx) msm.ResponseCheckTx {
	_, _, err := parseTx(req.Tx)
	if err != nil {
		return msm.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}
	return msm.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

// DeliverTx implements MSM.
func (app *Application) DeliverTx(req msm.RequestDeliverTx) msm.ResponseDeliverTx {
	key, value, err := parseTx(req.Tx)
	if err != nil {
		panic(err) // shouldn't happen since we verified it in CheckTx
	}
	app.state.Set(key, value)
	return msm.ResponseDeliverTx{Code: code.CodeTypeOK}
}

// EndBlock implements MSM.
func (app *Application) EndBlock(req msm.RequestEndBlock) msm.ResponseEndBlock {
	var err error
	resp := msm.ResponseEndBlock{}
	if resp.ValidatorUpdates, err = app.validatorUpdates(uint64(req.Height)); err != nil {
		panic(err)
	}
	return resp
}

// Commit implements MSM.
func (app *Application) Commit() msm.ResponseCommit {
	height, hash, err := app.state.Commit()
	if err != nil {
		panic(err)
	}
	if app.cfg.SnapshotInterval > 0 && height%app.cfg.SnapshotInterval == 0 {
		snapshot, err := app.snapshots.Create(app.state)
		if err != nil {
			panic(err)
		}
		logger.Info("Created state sync snapshot", "height", snapshot.Height)
	}
	retainHeight := int64(0)
	if app.cfg.RetainBlocks > 0 {
		retainHeight = int64(height - app.cfg.RetainBlocks + 1)
	}
	return msm.ResponseCommit{
		Data:         hash,
		RetainHeight: retainHeight,
	}
}

// Query implements MSM.
func (app *Application) Query(req msm.RequestQuery) msm.ResponseQuery {
	return msm.ResponseQuery{
		Height: int64(app.state.Height),
		Key:    req.Data,
		Value:  []byte(app.state.Get(string(req.Data))),
	}
}

// ListSnapshots implements MSM.
func (app *Application) ListSnapshots(req msm.RequestListSnapshots) msm.ResponseListSnapshots {
	snapshots, err := app.snapshots.List()
	if err != nil {
		panic(err)
	}
	return msm.ResponseListSnapshots{Snapshots: snapshots}
}

// LoadSnapshotChunk implements MSM.
func (app *Application) LoadSnapshotChunk(req msm.RequestLoadSnapshotChunk) msm.ResponseLoadSnapshotChunk {
	chunk, err := app.snapshots.LoadChunk(req.Height, req.Format, req.Chunk)
	if err != nil {
		panic(err)
	}
	return msm.ResponseLoadSnapshotChunk{Chunk: chunk}
}

// OfferSnapshot implements MSM.
func (app *Application) OfferSnapshot(req msm.RequestOfferSnapshot) msm.ResponseOfferSnapshot {
	if app.restoreSnapshot != nil {
		panic("A snapshot is already being restored")
	}
	app.restoreSnapshot = req.Snapshot
	app.restoreChunks = [][]byte{}
	return msm.ResponseOfferSnapshot{Result: msm.ResponseOfferSnapshot_ACCEPT}
}

// ApplySnapshotChunk implements MSM.
func (app *Application) ApplySnapshotChunk(req msm.RequestApplySnapshotChunk) msm.ResponseApplySnapshotChunk {
	if app.restoreSnapshot == nil {
		panic("No restore in progress")
	}
	app.restoreChunks = append(app.restoreChunks, req.Chunk)
	if len(app.restoreChunks) == int(app.restoreSnapshot.Chunks) {
		bz := []byte{}
		for _, chunk := range app.restoreChunks {
			bz = append(bz, chunk...)
		}
		err := app.state.Import(app.restoreSnapshot.Height, bz)
		if err != nil {
			panic(err)
		}
		app.restoreSnapshot = nil
		app.restoreChunks = nil
	}
	return msm.ResponseApplySnapshotChunk{Result: msm.ResponseApplySnapshotChunk_ACCEPT}
}

// validatorUpdates generates a validator set update.
func (app *Application) validatorUpdates(height uint64) (msm.ValidatorUpdates, error) {
	updates := app.cfg.ValidatorUpdates[fmt.Sprintf("%v", height)]
	if len(updates) == 0 {
		return nil, nil
	}

	valUpdates := msm.ValidatorUpdates{}
	for keyString, power := range updates {

		keyBytes, err := base64.StdEncoding.DecodeString(keyString)
		if err != nil {
			return nil, fmt.Errorf("invalid base64 pubkey value %q: %w", keyString, err)
		}
		valUpdates = append(valUpdates, msm.UpdateValidator(keyBytes, int64(power), app.cfg.KeyType))
	}
	return valUpdates, nil
}

// parseTx parses a tx in 'key=value' format into a key and value.
func parseTx(tx []byte) (string, string, error) {
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid tx format: %q", string(tx))
	}
	if len(parts[0]) == 0 {
		return "", "", errors.New("key cannot be empty")
	}
	return string(parts[0]), string(parts[1]), nil
}
