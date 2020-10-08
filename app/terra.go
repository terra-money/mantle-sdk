package app

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	TerraApp "github.com/terra-project/core/app"
	core "github.com/terra-project/core/types"
	compatapp "github.com/terra-project/mantle-compatibility/app"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
	l "log"
)

type App struct {
	terra *TerraApp.TerraApp
	db    dbm.DB
}

var (
	GlobalTerraApp *TerraApp.TerraApp
)

func NewApp(
	db dbm.DB,
	genesis *tmtypes.GenesisDoc,
) *App {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)
	config.Seal()

	app := compatapp.NewTerraApp(db)

	// only init chain on genesis
	if app.LastBlockHeight() == 0 {
		// init chain
		validators := make([]*tmtypes.Validator, len(genesis.Validators))
		for i, val := range genesis.Validators {
			validators[i] = tmtypes.NewValidator(val.PubKey, val.Power)
		}
		validatorSet := tmtypes.NewValidatorSet(validators)
		nextVals := tmtypes.TM2PB.ValidatorUpdates(validatorSet)
		csParams := tmtypes.TM2PB.ConsensusParams(genesis.ConsensusParams)
		ic := abci.RequestInitChain{
			Time:            genesis.GenesisTime,
			ChainId:         genesis.ChainID,
			AppStateBytes:   genesis.AppState,
			ConsensusParams: csParams,
			Validators:      nextVals,
		}

		initChainResponse := app.InitChain(ic)
		initChainResponseJSON, _ := json.Marshal(initChainResponse)

		validatorUpdate(db, 1, initChainResponse.Validators)
		consensusParamUpdate(db, 1, initChainResponse.ConsensusParams)

		commitResponse := app.Commit()
		commitResponseJSON, _ := json.Marshal(commitResponse)

		l.Printf("Init chain finished, LastBlockHeight=%d", app.LastBlockHeight())
		l.Printf("== InitChainResponse: %s", string(initChainResponseJSON))
		l.Printf("== CommitResponse: %s", string(commitResponseJSON))
	}

	GlobalTerraApp = app

	return &App{
		db:    db,
		terra: app,
	}
}

func (c *App) GetApp() *TerraApp.TerraApp {
	return c.terra
}

func (c *App) GetQueryRouter() sdk.QueryRouter {
	return c.terra.QueryRouter()
}

func (c *App) BeginBlocker(block *types.Block) abci.ResponseBeginBlock {
	var abciHeader = utils.ConvertToABCIHeader(&block.Header)

	// beginblock validator info
	commitInfo, byzVals := getBeginBlockValidatorInfo(block, c.db)

	var abciRequest = abci.RequestBeginBlock{
		Hash:                nil,
		Header:              abciHeader,
		LastCommitInfo:      commitInfo,
		ByzantineValidators: byzVals,
	}

	abciResponse := c.terra.BeginBlock(abciRequest)

	return abciResponse
}

func (c *App) EndBlocker(block *types.Block) abci.ResponseEndBlock {
	abciRequest := abci.RequestEndBlock{
		Height: block.Header.Height,
	}
	abciResponse := c.terra.EndBlock(abciRequest)

	// need to do validator, consensusParam update
	// only do this if nextHeight is NOT 1 (same as tendermint)
	defer func() {
		height := block.Header.Height
		validatorUpdate(c.db, height, abciResponse.ValidatorUpdates)
		consensusParamUpdate(c.db, height, abciResponse.ConsensusParamUpdates)
	}()

	return abciResponse
}

func (c *App) decodeTx(txbytes []byte) (sdk.Tx, error) {
	return types.TxDecoder(txbytes)
}

func (c *App) DeliverTxs(txs []tmtypes.Tx) []abci.ResponseDeliverTx {
	responses := make([]abci.ResponseDeliverTx, len(txs))
	for i, tx := range txs {
		response := c.terra.DeliverTx(abci.RequestDeliverTx{
			Tx: tx,
		})

		responses[i] = response
	}

	return responses
}

func (c *App) Commit(transactional bool) abci.ResponseCommit {
	response := c.terra.Commit()

	// no need to further save app state per block

	if c.terra.LastBlockHeight() == 192 {
		s, v, _ := c.terra.ExportAppStateAndValidators(false, []string{})
		fmt.Println(string(s), v)
	}

	if !transactional {
		return response
	}

	return response
}

func validatorUpdate(db dbm.DB, height int64, validatorUpdates abci.ValidatorUpdates) error {
	// may need to port this logic to mantle-compatibility,
	// as this logic tends to change with tendermint version
	var valInfo = &state.ValidatorsInfo{
		ValidatorSet:      nil,
		LastHeightChanged: height + 1,
	} //saveValidatorsInfo(db, nextHeight+1, state.LastHeightValidatorsChanged, state.NextValidators)

	//if height%100000 == 0 {
	validators, err := tmtypes.PB2TM.ValidatorUpdates(validatorUpdates)
	var validatorSet = tmtypes.NewValidatorSet(validators)

	if err != nil {
		return err
	}
	valInfo.ValidatorSet = validatorSet
	//}

	if len(valInfo.ValidatorSet.Validators) > 0 {
		fmt.Println("validator update")

	}

	return db.Set([]byte(fmt.Sprintf("validatorsKey:%v", height)), valInfo.Bytes())
}

func consensusParamUpdate(db dbm.DB, height int64, consensusParamUpdates *abci.ConsensusParams) {
	//tmtypes.
}

func getBeginBlockValidatorInfo(block *types.Block, stateDB dbm.DB) (abci.LastCommitInfo, []abci.Evidence) {
	voteInfos := make([]abci.VoteInfo, block.LastCommit.Size())
	// block.Height=1 -> LastCommitInfo.Votes are empty.
	// Remember that the first LastCommit is intentionally empty, so it makes
	// sense for LastCommitInfo.Votes to also be empty.
	if block.Header.Height > 1 {
		lastValSet, err := state.LoadValidators(stateDB, block.Header.Height-1)
		if err != nil {
			panic(err)
		}

		// Sanity check that commit size matches validator set size - only applies
		// after first block.
		var (
			commitSize = block.LastCommit.Size()
			valSetLen  = len(lastValSet.Validators)
		)
		if commitSize != valSetLen {
			panic(fmt.Sprintf("commit size (%d) doesn't match valset length (%d) at height %d\n\n%v\n\n%v",
				commitSize, valSetLen, block.Header.Height, block.LastCommit.Signatures, lastValSet.Validators))
		}

		for i, val := range lastValSet.Validators {
			commitSig := block.LastCommit.Signatures[i]
			voteInfos[i] = abci.VoteInfo{
				Validator:       tmtypes.TM2PB.Validator(val),
				SignedLastBlock: !commitSig.Absent(),
			}
		}
	}

	byzVals := make([]abci.Evidence, len(block.Evidence.Evidence))
	for i, ev := range block.Evidence.Evidence {
		// We need the validator set. We already did this in validateBlock.
		// TODO: Should we instead cache the valset in the evidence itself and add
		// `SetValidatorSet()` and `ToABCI` methods ?
		valset, err := state.LoadValidators(stateDB, ev.Height())
		if err != nil {
			panic(err)
		}
		byzVals[i] = tmtypes.TM2PB.Evidence(ev, valset, block.Header.Time)
	}

	// testkit blocks don't have LastCommit reference set,
	// treat this
	var round int32 = 0
	if block.LastCommit != nil {
		round = int32(block.LastCommit.Round)
	}

	return abci.LastCommitInfo{
		Round: round,
		Votes: voteInfos,
	}, byzVals
}
