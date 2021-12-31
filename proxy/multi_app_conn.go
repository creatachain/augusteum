package proxy

import (
	"fmt"

	tmlog "github.com/creatachain/augusteum/libs/log"
	tmos "github.com/creatachain/augusteum/libs/os"
	"github.com/creatachain/augusteum/libs/service"
	msmcli "github.com/creatachain/augusteum/msm/client"
)

const (
	connConsensus = "consensus"
	connMempool   = "mempool"
	connQuery     = "query"
	connSnapshot  = "snapshot"
)

// AppConns is the Augusteum's interface to the application that consists of
// multiple connections.
type AppConns interface {
	service.Service

	// Mempool connection
	Mempool() AppConnMempool
	// Consensus connection
	Consensus() AppConnConsensus
	// Query connection
	Query() AppConnQuery
	// Snapshot connection
	Snapshot() AppConnSnapshot
}

// NewAppConns calls NewMultiAppConn.
func NewAppConns(clientCreator ClientCreator) AppConns {
	return NewMultiAppConn(clientCreator)
}

// multiAppConn implements AppConns.
//
// A multiAppConn is made of a few appConns and manages their underlying msm
// clients.
// TODO: on app restart, clients must reboot together
type multiAppConn struct {
	service.BaseService

	consensusConn AppConnConsensus
	mempoolConn   AppConnMempool
	queryConn     AppConnQuery
	snapshotConn  AppConnSnapshot

	consensusConnClient msmcli.Client
	mempoolConnClient   msmcli.Client
	queryConnClient     msmcli.Client
	snapshotConnClient  msmcli.Client

	clientCreator ClientCreator
}

// NewMultiAppConn makes all necessary msm connections to the application.
func NewMultiAppConn(clientCreator ClientCreator) AppConns {
	multiAppConn := &multiAppConn{
		clientCreator: clientCreator,
	}
	multiAppConn.BaseService = *service.NewBaseService(nil, "multiAppConn", multiAppConn)
	return multiAppConn
}

func (app *multiAppConn) Mempool() AppConnMempool {
	return app.mempoolConn
}

func (app *multiAppConn) Consensus() AppConnConsensus {
	return app.consensusConn
}

func (app *multiAppConn) Query() AppConnQuery {
	return app.queryConn
}

func (app *multiAppConn) Snapshot() AppConnSnapshot {
	return app.snapshotConn
}

func (app *multiAppConn) OnStart() error {
	c, err := app.msmClientFor(connQuery)
	if err != nil {
		return err
	}
	app.queryConnClient = c
	app.queryConn = NewAppConnQuery(c)

	c, err = app.msmClientFor(connSnapshot)
	if err != nil {
		app.stopAllClients()
		return err
	}
	app.snapshotConnClient = c
	app.snapshotConn = NewAppConnSnapshot(c)

	c, err = app.msmClientFor(connMempool)
	if err != nil {
		app.stopAllClients()
		return err
	}
	app.mempoolConnClient = c
	app.mempoolConn = NewAppConnMempool(c)

	c, err = app.msmClientFor(connConsensus)
	if err != nil {
		app.stopAllClients()
		return err
	}
	app.consensusConnClient = c
	app.consensusConn = NewAppConnConsensus(c)

	// Kill Augusteum if the MSM application crashes.
	go app.killTMOnClientError()

	return nil
}

func (app *multiAppConn) OnStop() {
	app.stopAllClients()
}

func (app *multiAppConn) killTMOnClientError() {
	killFn := func(conn string, err error, logger tmlog.Logger) {
		logger.Error(
			fmt.Sprintf("%s connection terminated. Did the application crash? Please restart augusteum", conn),
			"err", err)
		killErr := tmos.Kill()
		if killErr != nil {
			logger.Error("Failed to kill this process - please do so manually", "err", killErr)
		}
	}

	select {
	case <-app.consensusConnClient.Quit():
		if err := app.consensusConnClient.Error(); err != nil {
			killFn(connConsensus, err, app.Logger)
		}
	case <-app.mempoolConnClient.Quit():
		if err := app.mempoolConnClient.Error(); err != nil {
			killFn(connMempool, err, app.Logger)
		}
	case <-app.queryConnClient.Quit():
		if err := app.queryConnClient.Error(); err != nil {
			killFn(connQuery, err, app.Logger)
		}
	case <-app.snapshotConnClient.Quit():
		if err := app.snapshotConnClient.Error(); err != nil {
			killFn(connSnapshot, err, app.Logger)
		}
	}
}

func (app *multiAppConn) stopAllClients() {
	if app.consensusConnClient != nil {
		if err := app.consensusConnClient.Stop(); err != nil {
			app.Logger.Error("error while stopping consensus client", "error", err)
		}
	}
	if app.mempoolConnClient != nil {
		if err := app.mempoolConnClient.Stop(); err != nil {
			app.Logger.Error("error while stopping mempool client", "error", err)
		}
	}
	if app.queryConnClient != nil {
		if err := app.queryConnClient.Stop(); err != nil {
			app.Logger.Error("error while stopping query client", "error", err)
		}
	}
	if app.snapshotConnClient != nil {
		if err := app.snapshotConnClient.Stop(); err != nil {
			app.Logger.Error("error while stopping snapshot client", "error", err)
		}
	}
}

func (app *multiAppConn) msmClientFor(conn string) (msmcli.Client, error) {
	c, err := app.clientCreator.NewMSMClient()
	if err != nil {
		return nil, fmt.Errorf("error creating MSM client (%s connection): %w", conn, err)
	}
	c.SetLogger(app.Logger.With("module", "msm-client", "connection", conn))
	if err := c.Start(); err != nil {
		return nil, fmt.Errorf("error starting MSM client (%s connection): %w", conn, err)
	}
	return c, nil
}
