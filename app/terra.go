package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	TerraApp "github.com/terra-project/core/app"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/auth"
	compatapp "github.com/terra-project/mantle-compatibility/app"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
	l "log"
)

type App struct {
	terra *TerraApp.TerraApp
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
		commitResponse := app.Commit()
		commitResponseJSON, _ := json.Marshal(commitResponse)

		l.Printf("Init chain finished, LastBlockHeight=%d", app.LastBlockHeight())
		l.Printf("== InitChainResponse: %s", string(initChainResponseJSON))
		l.Printf("== CommitResponse: %s", string(commitResponseJSON))
	}

	GlobalTerraApp = app

	return &App{
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
	var abciRequest = abci.RequestBeginBlock{
		Header: abciHeader,
	}

	abciResponse := c.terra.BeginBlock(abciRequest)

	return abciResponse
}

func (c *App) EndBlocker(block *types.Block) abci.ResponseEndBlock {
	abciRequest := abci.RequestEndBlock{
		Height: block.Header.Height,
	}
	abciResponse := c.terra.EndBlock(abciRequest)

	return abciResponse
}

func (c *App) decodeTx(txbytes []byte) (sdk.Tx, error) {
	var tx = auth.StdTx{}

	if len(txbytes) == 0 {
		return nil, fmt.Errorf("Tx bytes are empty")
	}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec
	err := c.terra.Codec().UnmarshalBinaryLengthPrefixed(txbytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (c *App) DeliverTxs(txs []string) []abci.ResponseDeliverTx {
	responses := make([]abci.ResponseDeliverTx, len(txs))
	for i, txstring := range txs {
		txbytes, err := base64.StdEncoding.DecodeString(txstring)
		if err != nil {
			panic(err)
		}

		response := c.terra.DeliverTx(abci.RequestDeliverTx{
			Tx: txbytes,
		})

		responses[i] = response
	}

	return responses
}

func (c *App) Commit(transactional bool) abci.ResponseCommit {
	response := c.terra.Commit()

	// no need to further save app state per block
	if !transactional {
		return response
	}
	// c.terra.ExportAppStateAndValidators()

	return response
}
