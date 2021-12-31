package version

var (
	// TMCoreSemVer is the current version of Augusteum Core.
	// It's the Semantic Version of the software.
	TMCoreSemVer string
)

const (
	// MSMSemVer is the semantic version of the MSM library
	MSMSemVer = "0.17.0"

	MSMVersion = MSMSemVer
)

var (
	// P2PProtocol versions all p2p behaviour and msgs.
	// This includes proposer selection.
	P2PProtocol uint64 = 8

	// BlockProtocol versions all block data structures and processing.
	// This includes validity of blocks and state updates.
	BlockProtocol uint64 = 11
)
