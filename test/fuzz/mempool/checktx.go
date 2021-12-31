package checktx

import (
	"github.com/creatachain/augusteum/config"
	mempl "github.com/creatachain/augusteum/mempool"
	"github.com/creatachain/augusteum/msm/example/kvstore"
	"github.com/creatachain/augusteum/proxy"
)

var mempool mempl.Mempool

func init() {
	app := kvstore.NewApplication()
	cc := proxy.NewLocalClientCreator(app)
	appConnMem, _ := cc.NewMSMClient()
	err := appConnMem.Start()
	if err != nil {
		panic(err)
	}

	cfg := config.DefaultMempoolConfig()
	cfg.Broadcast = false

	mempool = mempl.NewCListMempool(cfg, appConnMem, 0)
}

func Fuzz(data []byte) int {
	err := mempool.CheckTx(data, nil, mempl.TxInfo{})
	if err != nil {
		return 0
	}

	return 1
}
