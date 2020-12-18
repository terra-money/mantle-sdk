package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/auth"
)

type AutomaticTxEntry struct {
	Msgs        []sdk.Msg
	Fee         auth.StdFee
	Period      int
	AccountName string
}

func NewAutomaticTxEntry(accountName string, fee auth.StdFee, msgs []sdk.Msg, period int) AutomaticTxEntry {
	return AutomaticTxEntry{
		Msgs:        msgs,
		Fee:         fee,
		AccountName: accountName,
		Period:      period,
	}
}
