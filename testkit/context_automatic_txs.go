package testkit

import (
	"fmt"
	"log"
)

func (ctx *TestkitContext) GetAutomaticTxEntries() (atxs []AutomaticTxEntry) {
	return ctx.autoTxs
}

func (ctx *TestkitContext) AddAutomaticTxEntry(entry AutomaticTxEntry) (atxId string) {
	ctx.m.Lock()

	if entry.StartedAt == 0 {
		entry.StartedAt = ctx.mantle.GetApp().BaseApp.LastBlockHeight() + 1
	}

	log.Printf("[mantle/testkit/context_automatic_txs] registered atx %s\n", entry.ID)
	log.Printf("[mantle/testkit/context_automatic_txs]\tfrom accountName=%s, period=%d, startedAt=%d", entry.AccountName, entry.Period, entry.StartedAt)

	ctx.autoTxs = append(ctx.autoTxs, entry)
	ctx.m.Unlock()

	return entry.ID
}

func (ctx *TestkitContext) ClearAllAutomaticTxEntries() {
	ctx.m.Lock()
	ctx.autoTxs = make([]AutomaticTxEntry, 0)
	ctx.m.Unlock()
}

func (ctx *TestkitContext) ClearAutomaticTxEntry(atxId string) {
	ctx.m.Lock()
	defer ctx.m.Unlock()

	var isEntryFound = false
	nextAutoTxs := make([]AutomaticTxEntry, 0)

	for _, atxEntry := range ctx.autoTxs {
		if atxEntry.ID == atxId {
			isEntryFound = true
			continue
		}
		nextAutoTxs = append(nextAutoTxs, atxEntry)
	}

	if isEntryFound == false {
		panic(fmt.Errorf("atx entry %s is not found", atxId))
	}

	ctx.autoTxs = nextAutoTxs

	log.Printf("[mantle/testkit/context_automatic_txs] removed atx %s\n", atxId)
}
