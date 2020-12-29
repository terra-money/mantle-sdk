package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/terra-project/core/x/auth"
)

type AutomaticTxEntry struct {
	ID          string      `json:"id"`
	Msgs        []sdk.Msg   `json:"msgs"`
	Fee         auth.StdFee `json:"fee"`
	Period      int         `json:"period"`
	StartedAt   int64       `json:"started_at"`
	AccountName string      `json:"account_name"`
}

func NewAutomaticTxEntry(
	accountName string,
	fee auth.StdFee,
	msgs []sdk.Msg,
	period int,
	offset int64,
	startAt int64,
) AutomaticTxEntry {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return AutomaticTxEntry{
		ID:          uid.String(),
		Msgs:        msgs,
		Fee:         fee,
		AccountName: accountName,
		Period:      period,
		StartedAt:   startAt,
	}
}
