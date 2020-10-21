package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

// all types defined here MUST follow protobuf version of them.
//
// TODO: move this part to mantle-compatibility
type (
	Tx         = tmtypes.Tx
	GenesisDoc = tmtypes.GenesisDoc
	Block      = tmtypes.Block
	Header     = tmtypes.Header
)

//
//type Block struct {
//	Header     tmtypes.Header       `json:"header"`
//	Data       tmtypes.Data         `json:"data"`
//	Evidence   tmtypes.EvidenceData `json:"evidence"`
//	LastCommit *tmtypes.Commit      `json:"last_commit"`
//}
//
//// App includes the protocol and software version for the application.
//// This information is included in ResponseInfo. The App.Protocol can be
//// updated in ResponseEndBlock.
//type App struct {
//	Protocol version.Protocol `json:"protocol,string"`
//	Software string           `json:"software"`
//}
//
//// Consensus captures the consensus rules for processing a block in the blockchain,
//// including all blockchain data structures and the rules of the application's
//// state transition machine.
//type Consensus struct {
//	Block version.Protocol `json:"block,string"`
//	App   version.Protocol `json:"app,string"`
//}
//
//type BlockID struct {
//	Hash        compattypes.HexBytes `json:"hash"`
//	PartsHeader PartSetHeader        `json:"parts"`
//}
//
//type PartSetHeader struct {
//	Total int                  `json:"total,string"`
//	Hash  compattypes.HexBytes `json:"hash"`
//}
//
////
////type Header struct {
////	// basic block info
////	Version Consensus `json:"version"`
////	ChainID string    `json:"chain_id"`
////	Height  int64     `json:"height,string"`
////	Time    time.Time `json:"time"`
////
////	// prev block info
////	LastBlockID BlockID `json:"last_block_id"`
////
////	// hashes of block data
////	LastCommitHash compattypes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
////	DataHash       compattypes.HexBytes `json:"data_hash"`        // transactions
////
////	// hashes from the app output from the prev block
////	ValidatorsHash     compattypes.HexBytes `json:"validators_hash"`      // validators for the current block
////	NextValidatorsHash compattypes.HexBytes `json:"next_validators_hash"` // validators for the next block
////	ConsensusHash      compattypes.HexBytes `json:"consensus_hash"`       // consensus params for current block
////	AppHash            compattypes.HexBytes `json:"app_hash"`             // state after txs from the previous block
////	// root hash of all results from the txs from the previous block
////	LastResultsHash compattypes.HexBytes `json:"last_results_hash"`
////
////	// consensus info
////	EvidenceHash    compattypes.HexBytes `json:"evidence_hash"`    // evidence included in the block
////	ProposerAddress tmtypes.Address      `json:"proposer_address"` // original proposer of the block
////}
////
////type HeaderCommit struct {
////	// NOTE: The signatures are in order of address to preserve the bonded
////	// ValidatorSet order.
////	// Any peer with a block can gossip signatures by index with a peer without
////	// recalculating the active ValidatorSet.
////	Height     int64               `json:"height,string"`
////	Round      int                 `json:"round,string"`
////	BlockID    BlockID             `json:"block_id"`
////	Signatures []tmtypes.CommitSig `json:"signatures"`
////
////	// Memoized in first call to corresponding method.
////	// NOTE: can't memoize in constructor because constructor isn't used for
////	// unmarshaling.
////	hash     compattypes.HexBytes
////	bitArray *bits.BitArray
////}
//
////
////type Header struct {
////	Version struct {
////		Block uint64 `json:"Block,string"`
////		App   uint64 `json:"App,string"`
////	} `json:"version"`
////	ChainID string `json:"chain_id"`
////	Height  int64  `json:"height,string"`
////	Time    string `json:"time"`
////	// NumTxs      int64  `json:"num_txs,string"` // removed as of mantlemint 0.33
////	// TotalTxs    int64  `json:"total_txs,string"` // remove as of mantlemint 0.33
////	LastBlockID struct {
////		Hash  string `json:"hash"`
////		Parts struct {
////			Total string `json:"total"`
////			Hash  string `json:"hash"`
////		} `json:"parts"`
////	} `json:"last_block_id"`
////	LastCommitHash     string `json:"last_commit_hash"`
////	DataHash           string `json:"data_hash"`
////	ValidatorsHash     string `json:"validators_hash"`
////	NextValidatorsHash string `json:"next_validators_hash"`
////	ConsensusHash      string `json:"consensus_hash"`
////	AppHash            string `json:"app_hash"`
////	LastResultsHash    string `json:"last_results_hash"`
////	EvidenceHash       string `json:"evidence_hash"`
////	ProposerAddress    string `json:"proposer_address"`
////}
////
////type ResultBeginBlock struct {
////	Events []struct {
////		Type       string `json:"type"`
////		Attributes []struct {
////			Key   string `json:"key"`
////			Value string `json:"value"`
////		} `json:"attributes"`
////	} `json:"events"`
////}
//
//type ResultEndBlock struct {
//	ValidatorUpdates []interface{} `json:"validator_updates"`
//	Events           []struct {
//		Type       string `json:"type"`
//		Attributes []struct {
//			Key   string `json:"key"`
//			Value string `json:"value"`
//		} `json:"attributes"`
//	} `json:"events"`
//}
//
//type Events struct {
//	RewardsAmount                  []string `json:"rewards.amount"`
//	RewardsValidator               []string `json:"rewards.validator"`
//	ExchangeRateUpdateExchangeRate []string `json:"exchange_rate_update.exchange_rate"`
//	TransferRecipient              []string `json:"transfer.recipient"`
//	TransferAmount                 []string `json:"transfer.amount"`
//	MessageSender                  []string `json:"message.sender"`
//	ProposerRewardValidator        []string `json:"proposer_reward.validator"`
//	CommissionValidator            []string `json:"commission.validator"`
//	ProposerRewardAmount           []string `json:"proposer_reward.amount"`
//	CommissionAmount               []string `json:"commission.amount"`
//	ExchangeRateUpdateDenom        []string `json:"exchange_rate_update.denom"`
//	TmEvent                        []string `json:"tm.event"`
//}
