package testkit

func (ctx *TestkitContext) AddAutomaticTxEntry(entry AutomaticTxEntry) {
	ctx.m.Lock()
	ctx.autoTxs = append(ctx.autoTxs, entry)
	ctx.m.Unlock()
}
