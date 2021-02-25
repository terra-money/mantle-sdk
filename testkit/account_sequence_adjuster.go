package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	terra "github.com/terra-project/core/app"
	"github.com/terra-project/core/x/auth"
)

type AccountSequenceAdjuster struct {
	isPurged      bool
	querier       auth.NodeQuerier
	ar            auth.AccountRetriever
	app           *terra.TerraApp
	adjustmentMap map[string]exported.Account
}

func NewAccountSequenceAdjuster(app *terra.TerraApp) *AccountSequenceAdjuster {
	querier := NewLocalQuerier(app)
	return &AccountSequenceAdjuster{
		querier:       querier,
		app:           app,
		ar:            auth.NewAccountRetriever(querier),
		adjustmentMap: make(map[string]exported.Account),
	}
}

func (asa *AccountSequenceAdjuster) GetNextSequence(address sdk.AccAddress) (account exported.Account) {
	if asa.isPurged {
		panic("cannot reuse purged account sequence adjuster")
	}

	lastKnownAccount, ok := asa.adjustmentMap[address.String()]
	if !ok {
		account := asa.GetOrCreateAccount(address)
		asa.adjustmentMap[address.String()] = account
		lastKnownAccount = account
	} else {
		if err := lastKnownAccount.SetSequence(lastKnownAccount.GetSequence() + 1); err != nil {
			panic(err)
		}
	}

	return lastKnownAccount
}

func (asa *AccountSequenceAdjuster) Purge() {
	asa.isPurged = true
}

func (asa *AccountSequenceAdjuster) GetOrCreateAccount(address sdk.AccAddress) exported.Account {
	ar := auth.NewAccountRetriever(asa.querier)

	// check if account is known
	if accountExists := ar.EnsureExists(address); accountExists != nil {
		panic(accountExists)
	}

	acc, accErr := ar.GetAccount(address)
	if accErr != nil {
		panic(accErr)
	}

	return acc
}
