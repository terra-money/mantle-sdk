package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/terra-project/core/x/gov"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/oracle"
	"github.com/terra-project/core/x/slashing"
	"github.com/terra-project/core/x/staking"
	"github.com/terra-project/core/x/wasm"
)

// bank modules are set as internal,
// so redeclare them here
func NewMsgSend(
	fromAddress sdk.AccAddress,
	toAddress sdk.AccAddress,
	amount sdk.Coins,
) bank.MsgSend {
	return bank.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}
}

func NewMsgMultiSend(
	Inputs []bank.Input,
	Outputs []bank.Output,
) bank.MsgMultiSend {
	return bank.MsgMultiSend{
		Inputs:  Inputs,
		Outputs: Outputs,
	}
}

// distribution
type (
	MsgModifyWithdrawAddress       = distribution.MsgSetWithdrawAddress
	MsgWithdrawDelegationReward    = distribution.MsgWithdrawDelegatorReward
	MsgWithdrawValidatorCommission = distribution.MsgWithdrawValidatorCommission
)

var (
	NewMsgModifyWithdrawAddress       = distribution.NewMsgSetWithdrawAddress
	NewMsgWithdrawDelegationReward    = distribution.NewMsgWithdrawDelegatorReward
	NewMsgWithdrawValidatorCommission = distribution.NewMsgWithdrawValidatorCommission
	// NewFund

)

// governance
type (
	MsgDeposit        = gov.MsgDeposit
	MsgSubmitProposal = gov.MsgSubmitProposal
	MsgVote           = gov.MsgVote
)

var (
	NewMsgDeposit        = gov.NewMsgDeposit
	NewMsgSubmitProposal = gov.NewMsgSubmitProposal
	NewMsgVote           = gov.NewMsgVote
)

// market
type (
	MsgSwap = market.MsgSwap
	// MsgSwapSend = market.MsgSwapSend
)

var (
	NewMsgSwap = market.NewMsgSwap
)

// oracle
type (
	MsgExchangeRateVote    = oracle.MsgExchangeRateVote
	MsgExchangeRatePrevote = oracle.MsgExchangeRatePrevote
	MsgDelegateFeedConsent = oracle.MsgDelegateFeedConsent
)

var (
	NewMsgExchangeRateVote    = oracle.NewMsgExchangeRateVote
	NewMsgExchangeRatePrevote = oracle.NewMsgExchangeRatePrevote
	NewMsgDelegateFeedConsent = oracle.NewMsgDelegateFeedConsent
	// case 'oracle/MsgAggregateExchangeRatePrevote':
	// return MsgAggregateExchangeRatePrevote.fromData(data);
	// case 'oracle/MsgAggregateExchangeRateVote':
	// return MsgAggregateExchangeRateVote.fromData(data);
)

// slashing
type (
	MsgUnjail = slashing.MsgUnjail
)

var (
	NewMsgUnjail = slashing.NewMsgUnjail
)

// staking
type (
	MsgDelegate        = staking.MsgDelegate
	MsgUndelegate      = staking.MsgUndelegate
	MsgBeginRedelegate = staking.MsgBeginRedelegate
	MsgCreateValidator = staking.MsgCreateValidator
	MsgEditValidator   = staking.MsgEditValidator
)

var (
	NewMsgDelegate        = staking.NewMsgDelegate
	NewMsgUndelegate      = staking.NewMsgUndelegate
	NewMsgBeginRedelegate = staking.NewMsgBeginRedelegate
	NewMsgCreateValidator = staking.NewMsgCreateValidator
	NewMsgEditValidator   = staking.NewMsgEditValidator
)

//
// // wasm
// case 'wasm/MsgStoreCode':
// return MsgStoreCode.fromData(data);
// case 'wasm/MsgInstantiateContract':
// return MsgInstantiateContract.fromData(data);
// case 'wasm/MsgExecuteContract':
// return MsgExecuteContract.fromData(data);
// case 'wasm/MsgMigrateContract':
// return MsgMigrateContract.fromData(data);
// case 'wasm/MsgUpdateContractOwner':
// return MsgUpdateContractOwner.fromData(data);
type (
	MsgStoreCode           = wasm.MsgStoreCode
	MsgInstantiateContract = wasm.MsgInstantiateContract
	MsgExecuteContract     = wasm.MsgExecuteContract
	MsgMigrateContract     = wasm.MsgMigrateContract
	MsgUpdateContractOwner = wasm.MsgUpdateContractOwner
)

var (
	NewMsgStoreCode           = wasm.NewMsgStoreCode
	NewMsgInstantiateContract = wasm.NewMsgInstantiateContract
	NewMsgExecuteContract     = wasm.NewMsgExecuteContract
	NewMigrateContract        = wasm.NewMsgMigrateContract
	NewMsgUpdateContractOwner = wasm.NewMsgUpdateContractOwner
)
