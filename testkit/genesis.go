package testkit

import (
	"encoding/json"
	"fmt"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cauth "github.com/cosmos/cosmos-sdk/x/auth/exported"
	cauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/app"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/core/x/genutil"
	"time"
)

var defaultPassphrase = "12345678"

type TestkitAccount struct {
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic"`
}

type TestkitGenesisAccountToPrivValMap struct {
	Account       sdk.ValAddress
	PrivValidator tm.PrivValidator
}

type TestkitGenesis struct {
	chainId           string
	genesis           *tm.GenesisDoc
	isSealed          bool
	accounts          []TestkitAccount
	genesisValidators []auth.StdTx
	validatorMap      []TestkitGenesisAccountToPrivValMap
	kb                keys.Keybase

	hdAccount uint64
}

var ZeroCommission = types.NewCommissionRates(sdk.NewDec(0), sdk.NewDec(1), sdk.NewDec(1))

func NewTestkitGenesis(chainId string) *TestkitGenesis {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)

	return &TestkitGenesis{
		chainId:           chainId,
		genesis:           nil,
		isSealed:          false,
		accounts:          make([]TestkitAccount, 0),
		genesisValidators: []auth.StdTx{},
		validatorMap:      []TestkitGenesisAccountToPrivValMap{},
		kb:                ckeys.NewInMemoryKeyBase(),
		hdAccount:         0,
	}
}

func (tg *TestkitGenesis) IsSealed() bool {
	return tg.isSealed
}

// AddAccount creates new account from aan account name
func (tg *TestkitGenesis) AddAccount(accountName string) (keys.Info, string) {
	if tg.isSealed {
		panic("genesis sealed")
	}

	_, err := tg.kb.Get(accountName)
	if err == nil {
		panic(fmt.Errorf("can't re-add account %s", accountName))
	}

	i, seed, err := tg.kb.CreateMnemonic(
		accountName,
		keys.English,
		defaultPassphrase,
		keys.Secp256k1,
	)

	if err != nil {
		panic(err)
	}

	tg.accounts = append(tg.accounts, TestkitAccount{
		Address:  i.GetAddress().String(),
		Mnemonic: seed,
	})

	return i, seed
}

// ImportAccount imports account from provided account name and mnemonic
func (tg *TestkitGenesis) ImportAccount(accountName string, seed string) (keys.Info, string) {
	if tg.isSealed {
		panic("genesis sealed")
	}

	info, err := tg.kb.CreateAccount(
		accountName,
		seed,
		keys.DefaultBIP39Passphrase,
		defaultPassphrase,
		keys.CreateHDPath(0, 0).String(),
		keys.Secp256k1,
	)
	if err != nil {
		panic(err)
	}

	tg.accounts = append(tg.accounts, TestkitAccount{
		Address:  info.GetAddress().String(),
		Mnemonic: seed,
	})

	return info, seed
}

func (tg *TestkitGenesis) CreateValidator(
	accountName string,
	selfDelegation sdk.Coin,
	commission types.CommissionRates,
) types.MsgCreateValidator {
	if tg.isSealed {
		panic("genesis sealed")
	}
	key, err := tg.kb.Get(accountName)
	if err != nil {
		panic(err)
	}

	// create random privKey;
	// note that it doesn't really matter here
	pv := tm.NewMockPV()
	pvPub, pvPubErr := pv.GetPubKey()
	if pvPubErr != nil {
		panic(pvPubErr)
	}

	// convert to validator address
	valAddress := sdk.ValAddress(key.GetAddress().Bytes())

	// create validator message
	createValidatorMessage := types.NewMsgCreateValidator(
		valAddress,
		pvPub,
		selfDelegation,
		types.NewDescription(
			fmt.Sprintf("Validator (%s)", accountName),
			accountName,
			"",
			"",
			"",
		),
		commission,
		sdk.NewInt(1),
	)

	tx := auth.NewStdTx(
		[]sdk.Msg{createValidatorMessage},
		auth.NewStdFee(200000, nil),
		nil,
		"",
	)

	stdSignature, err := auth.MakeSignature(
		tg.kb,
		accountName,
		defaultPassphrase,
		auth.StdSignMsg{
			ChainID:       tg.chainId,
			AccountNumber: 0,
			Sequence:      0,
			Fee:           tx.Fee,
			Msgs:          tx.GetMsgs(),
			Memo:          tx.GetMemo(),
		},
	)

	tx.Signatures = []auth.StdSignature{stdSignature}

	// assign to genesis vals
	tg.genesisValidators = append(tg.genesisValidators, tx)
	tg.validatorMap = append(tg.validatorMap, TestkitGenesisAccountToPrivValMap{
		Account:       valAddress,
		PrivValidator: pv,
	})

	if err != nil {
		panic(err)
	}

	return createValidatorMessage
}

func (tg *TestkitGenesis) Seal() *tm.GenesisDoc {
	if tg.isSealed {
		panic("genesis sealed")
	}

	// create default genesis if TestkitGenesis is not created from
	// NewGenesisFromFile.
	// else reuse appstate from existing genesis
	var appState map[string]json.RawMessage
	if tg.genesis == nil {
		appState = app.ModuleBasics.DefaultGenesis()
	} else {
		tempAppState := make(map[string]json.RawMessage)
		json.Unmarshal(tg.genesis.AppState, &tempAppState)
		appState = tempAppState
	}

	// turn accounts into genesis accounts
	genesisAccounts := make([]cauth.GenesisAccount, len(tg.accounts))
	for gi, ga := range tg.accounts {
		addr, _ := sdk.AccAddressFromBech32(ga.Address)
		acc := cauthtypes.NewBaseAccountWithAddress(addr)
		acc.Coins = sdk.NewCoins(
			sdk.NewCoin("uluna", sdk.NewInt(100000000000000)),
			sdk.NewCoin("uusd", sdk.NewInt(100000000000000)),
		)

		genesisAccounts[gi] = &acc
	}

	authDefaultState := cauthtypes.DefaultGenesisState()
	authDefaultState.Accounts = genesisAccounts
	appState["auth"] = codec.MustMarshalJSON(authDefaultState)

	// genutil
	genutilDefaultState := genutil.NewGenesisStateFromStdTx(tg.genesisValidators)
	appState["genutil"] = codec.MustMarshalJSON(genutilDefaultState)

	// json marshal the whole thing
	genesisState, err := codec.MarshalJSON(appState)

	if err != nil {
		panic(err)
	}

	// never chainging ever since
	tg.isSealed = true

	gendoc := &tm.GenesisDoc{
		ChainID:     tg.chainId,
		Validators:  nil,
		AppState:    genesisState,
		GenesisTime: time.Now(),
	}

	if gendocErr := gendoc.ValidateAndComplete(); gendocErr != nil {
		panic(gendocErr)
	}

	tg.genesis = gendoc

	return gendoc
}

func (tg *TestkitGenesis) GetGenesisDoc() *tm.GenesisDoc {
	if !tg.isSealed {
		panic("genesis cannot be fetched before seal")
	}
	return tg.genesis
}

func (tg *TestkitGenesis) GetAccounts() []TestkitAccount {
	return tg.accounts
}

func (tg *TestkitGenesis) GetValidators() []TestkitGenesisAccountToPrivValMap {
	return tg.validatorMap
}

func (tg *TestkitGenesis) GetValidatorAddrs() []sdk.ValAddress {
	addrs := make([]sdk.ValAddress, len(tg.validatorMap))
	for i, v := range tg.validatorMap {
		addrs[i] = v.Account
	}

	return addrs
}

func (tg *TestkitGenesis) GetKeybase() keys.Keybase {
	return tg.kb
}
