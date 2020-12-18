package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AutomaticInjection struct {
	isEnabled    bool
	valRounds    []sdk.ValAddress
	lastProposer uint64
}

func (ctx *TestkitContext) SetAutomaticInjection(validatorRounds []string) {
	valAddresses := make([]sdk.ValAddress, len(validatorRounds))
	for i, accountName := range validatorRounds {
		info, err := ctx.tg.kb.Get(accountName)
		if err != nil {
			panic(err)
		}

		valAddresses[i] = sdk.ValAddress(info.GetAddress())
	}

	ctx.autoInjection.isEnabled = true
	ctx.autoInjection.valRounds = valAddresses
	ctx.autoInjection.lastProposer = 0
}

func (ai *AutomaticInjection) NextProposer() sdk.ValAddress {
	defer func() {
		ai.lastProposer = (ai.lastProposer + 1) % uint64(len(ai.valRounds))
	}()

	return ai.valRounds[ai.lastProposer]
}
