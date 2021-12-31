package client_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/creatachain/augusteum/msm/example/kvstore"
	nm "github.com/creatachain/augusteum/node"
	rpctest "github.com/creatachain/augusteum/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a augusteum node (and kvstore) in the background to test against
	dir, err := ioutil.TempDir("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}

	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartAugusteum(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopAugusteum(node)
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
