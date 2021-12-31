package proxy

import (
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/version"
)

// RequestInfo contains all the information for sending
// the msm.RequestInfo message during handshake with the app.
// It contains only compile-time version information.
var RequestInfo = msm.RequestInfo{
	Version:      version.TMCoreSemVer,
	BlockVersion: version.BlockProtocol,
	P2PVersion:   version.P2PProtocol,
}
