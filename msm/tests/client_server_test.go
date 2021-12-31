package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	msmclient "github.com/creatachain/augusteum/msm/client"
	"github.com/creatachain/augusteum/msm/example/kvstore"
	msmserver "github.com/creatachain/augusteum/msm/server"
)

func TestClientServerNoAddrPrefix(t *testing.T) {
	addr := "localhost:26658"
	transport := "socket"
	app := kvstore.NewApplication()

	server, err := msmserver.NewServer(addr, transport, app)
	assert.NoError(t, err, "expected no error on NewServer")
	err = server.Start()
	assert.NoError(t, err, "expected no error on server.Start")

	client, err := msmclient.NewClient(addr, transport, true)
	assert.NoError(t, err, "expected no error on NewClient")
	err = client.Start()
	assert.NoError(t, err, "expected no error on client.Start")
}
