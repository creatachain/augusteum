package types

import (
	"fmt"

	tmjson "github.com/creatachain/augusteum/libs/json"
	tmpubsub "github.com/creatachain/augusteum/libs/pubsub"
	tmquery "github.com/creatachain/augusteum/libs/pubsub/query"
	msm "github.com/creatachain/augusteum/msm/types"
)

// Reserved event types (alphabetically sorted).
const (
	// Block level events for mass consumption by users.
	// These events are triggered from the state package,
	// after a block has been committed.
	// These are also used by the tx indexer for async indexing.
	// All of this data can be fetched through the rpc.
	EventNewBlock            = "NewBlock"
	EventNewBlockHeader      = "NewBlockHeader"
	EventNewEvidence         = "NewEvidence"
	EventTx                  = "Tx"
	EventValidatorSetUpdates = "ValidatorSetUpdates"

	// Internal consensus events.
	// These are used for testing the consensus state machine.
	// They can also be used to build real-time consensus visualizers.
	EventCompleteProposal = "CompleteProposal"
	EventLock             = "Lock"
	EventNewRound         = "NewRound"
	EventNewRoundStep     = "NewRoundStep"
	EventPolka            = "Polka"
	EventRelock           = "Relock"
	EventTimeoutPropose   = "TimeoutPropose"
	EventTimeoutWait      = "TimeoutWait"
	EventUnlock           = "Unlock"
	EventValidBlock       = "ValidBlock"
	EventVote             = "Vote"
)

// ENCODING / DECODING

// TMEventData implements events.EventData.
type TMEventData interface {
	// empty interface
}

func init() {
	tmjson.RegisterType(EventDataNewBlock{}, "augusteum/event/NewBlock")
	tmjson.RegisterType(EventDataNewBlockHeader{}, "augusteum/event/NewBlockHeader")
	tmjson.RegisterType(EventDataNewEvidence{}, "augusteum/event/NewEvidence")
	tmjson.RegisterType(EventDataTx{}, "augusteum/event/Tx")
	tmjson.RegisterType(EventDataRoundState{}, "augusteum/event/RoundState")
	tmjson.RegisterType(EventDataNewRound{}, "augusteum/event/NewRound")
	tmjson.RegisterType(EventDataCompleteProposal{}, "augusteum/event/CompleteProposal")
	tmjson.RegisterType(EventDataVote{}, "augusteum/event/Vote")
	tmjson.RegisterType(EventDataValidatorSetUpdates{}, "augusteum/event/ValidatorSetUpdates")
	tmjson.RegisterType(EventDataString(""), "augusteum/event/ProposalString")
}

// Most event messages are basic types (a block, a transaction)
// but some (an input to a call tx or a receive) are more exotic

type EventDataNewBlock struct {
	Block *Block `json:"block"`

	ResultBeginBlock msm.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   msm.ResponseEndBlock   `json:"result_end_block"`
}

type EventDataNewBlockHeader struct {
	Header Header `json:"header"`

	NumTxs           int64                  `json:"num_txs"` // Number of txs in a block
	ResultBeginBlock msm.ResponseBeginBlock `json:"result_begin_block"`
	ResultEndBlock   msm.ResponseEndBlock   `json:"result_end_block"`
}

type EventDataNewEvidence struct {
	Evidence Evidence `json:"evidence"`

	Height int64 `json:"height"`
}

// All txs fire EventDataTx
type EventDataTx struct {
	msm.TxResult
}

// NOTE: This goes into the replay WAL
type EventDataRoundState struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`
}

type ValidatorInfo struct {
	Address Address `json:"address"`
	Index   int32   `json:"index"`
}

type EventDataNewRound struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`

	Proposer ValidatorInfo `json:"proposer"`
}

type EventDataCompleteProposal struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`

	BlockID BlockID `json:"block_id"`
}

type EventDataVote struct {
	Vote *Vote
}

type EventDataString string

type EventDataValidatorSetUpdates struct {
	ValidatorUpdates []*Validator `json:"validator_updates"`
}

// PUBSUB

const (
	// EventTypeKey is a reserved composite key for event name.
	EventTypeKey = "tm.event"
	// TxHashKey is a reserved key, used to specify transaction's hash.
	// see EventBus#PublishEventTx
	TxHashKey = "tx.hash"
	// TxHeightKey is a reserved key, used to specify transaction block's height.
	// see EventBus#PublishEventTx
	TxHeightKey = "tx.height"
)

var (
	EventQueryCompleteProposal    = QueryForEvent(EventCompleteProposal)
	EventQueryLock                = QueryForEvent(EventLock)
	EventQueryNewBlock            = QueryForEvent(EventNewBlock)
	EventQueryNewBlockHeader      = QueryForEvent(EventNewBlockHeader)
	EventQueryNewEvidence         = QueryForEvent(EventNewEvidence)
	EventQueryNewRound            = QueryForEvent(EventNewRound)
	EventQueryNewRoundStep        = QueryForEvent(EventNewRoundStep)
	EventQueryPolka               = QueryForEvent(EventPolka)
	EventQueryRelock              = QueryForEvent(EventRelock)
	EventQueryTimeoutPropose      = QueryForEvent(EventTimeoutPropose)
	EventQueryTimeoutWait         = QueryForEvent(EventTimeoutWait)
	EventQueryTx                  = QueryForEvent(EventTx)
	EventQueryUnlock              = QueryForEvent(EventUnlock)
	EventQueryValidatorSetUpdates = QueryForEvent(EventValidatorSetUpdates)
	EventQueryValidBlock          = QueryForEvent(EventValidBlock)
	EventQueryVote                = QueryForEvent(EventVote)
)

func EventQueryTxFor(tx Tx) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s='%X'", EventTypeKey, EventTx, TxHashKey, tx.Hash()))
}

func QueryForEvent(eventType string) tmpubsub.Query {
	return tmquery.MustParse(fmt.Sprintf("%s='%s'", EventTypeKey, eventType))
}

// BlockEventPublisher publishes all block related events
type BlockEventPublisher interface {
	PublishEventNewBlock(block EventDataNewBlock) error
	PublishEventNewBlockHeader(header EventDataNewBlockHeader) error
	PublishEventNewEvidence(evidence EventDataNewEvidence) error
	PublishEventTx(EventDataTx) error
	PublishEventValidatorSetUpdates(EventDataValidatorSetUpdates) error
}

type TxEventPublisher interface {
	PublishEventTx(EventDataTx) error
}
